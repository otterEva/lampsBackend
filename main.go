package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/otterEva/lamps/app/docs"
	"github.com/otterEva/lamps/app/handlers"
	"github.com/otterEva/lamps/app/middlewares"
	"github.com/otterEva/lamps/app/settings"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "lamps API",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
		_ = app.Shutdown()
	}()

	// -----------------------------------------------------------------

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// app.Get("/images/:image_url", handlers.GetImageHandler)
	app.Get("/images/:image_url", func(c *fiber.Ctx) error {
		log.Debug("Requested:", c.Params("image_url"))
		return handlers.GetImageHandler(c)
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	handlers.AuthHandlers(app.Group("/auth"), settings.Clients.DbClient, ctx)
	handlers.UserGoodsHandlers(app.Group("/goods"), settings.Clients.DbClient, ctx)

	protected := app.Group("/", middlewares.AuthMiddleware(settings.Clients.DbClient, ctx))

	handlers.UserOrdersHandler(protected.Group("/orders"), settings.Clients.DbClient, ctx)
	handlers.AdminGoodsHandler(protected.Group("/admin/goods"), settings.Clients.DbClient, ctx)
	handlers.AdminOrdersHandler(protected.Group("/admin/orders"), settings.Clients.DbClient, ctx)

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

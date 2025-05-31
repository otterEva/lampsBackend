package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/otterEva/lamps/goods_service/handlers"
	"github.com/otterEva/lamps/goods_service/middlewares"
	"github.com/otterEva/lamps/goods_service/settings"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://127.0.0.1:5173",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	handlers.UserGoodsHandlers(app.Group("/goods"), settings.Clients.DbClient, ctx)
	protected := app.Group("/", middlewares.AuthMiddleware(settings.Clients.DbClient, ctx))
	handlers.AdminGoodsHandler(protected.Group("/admin/goods"), settings.Clients.DbClient, ctx)

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

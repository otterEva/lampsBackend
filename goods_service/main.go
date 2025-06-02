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

	_ = ctx

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
		_ = app.Shutdown()
	}()

	// -----------------------------------------------------------------

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/goods", func(c *fiber.Ctx) error {
		return handlers.UserGoodsGet(c, ctx)
	})

	app.Get("/goods/:id", func(c *fiber.Ctx) error {
		return handlers.CheckIfGoodExists(c, ctx)
	})

	protected := app.Group("/", middlewares.AuthMiddleware())

	protected.Delete("/goods/admin/:id", func(c *fiber.Ctx) error {
		return handlers.AdminGoodDelete(c, ctx)
	})

	protected.Post("/goods/admin", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsPost(c, ctx)
	})

	protected.Get("/goods/admin", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsGet(c, ctx)
	})

	protected.Patch("/goods/admin/:id", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsPatch(c, ctx)
	})

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

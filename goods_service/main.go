package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
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

	app.Get("/goods", func(c *fiber.Ctx) error {
		return handlers.UserGoodsGet(c, ctx)
	})

	protected := app.Group("/goods/admin", middlewares.AuthMiddleware())

	protected.Post("", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsPost(c, ctx)
	})

	protected.Get("", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsGet(c, ctx)
	})

	protected.Delete("/:id", func(c *fiber.Ctx) error {
		return handlers.AdminGoodDelete(c, ctx)
	})

	protected.Patch("/:id", func(c *fiber.Ctx) error {
		return handlers.AdminGoodsPatch(c, ctx)
	})

	app.Get("/goods/:id", func(c *fiber.Ctx) error {
		return handlers.CheckIfGoodExists(c, ctx)
	})

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/orders_service/handlers"
	"github.com/otterEva/lamps/orders_service/middlewares"
	"github.com/otterEva/lamps/orders_service/settings"
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

	protected := app.Group("/", middlewares.AuthMiddleware(ctx))

	protected.Get("/orders", func(c *fiber.Ctx) error {
		return handlers.UserGetOrders(c, ctx)
	})

	protected.Post("/orders", func(c *fiber.Ctx) error {
		return handlers.UserPostOrder(c, ctx)
	})

	protected.Get("/orders/admin", func(c *fiber.Ctx) error {
		return handlers.AdminGetOrders(c, ctx)
	})

	protected.Delete("/orders/admin", func(c *fiber.Ctx) error {
		return handlers.AdminDeleteOrder(c, ctx)
	})

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

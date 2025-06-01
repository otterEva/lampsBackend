package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://127.0.0.1:5173",
	}))

	protected := app.Group("/", middlewares.AuthMiddleware(ctx))

	protected.Get("/orders", func (c *fiber.Ctx) error {
		return handlers.UserGetOrders(c, ctx)
	})
	protected.Get("/admin/orders", func (c *fiber.Ctx) error {
		return handlers.AdminGetOrders(c, ctx)
	})
	protected.Post("orders", func (c *fiber.Ctx) error {
		return handlers.UserPostOrder(c, ctx)
	})
	protected.Delete("/admin/orders", func (c *fiber.Ctx) error {
		return handlers.AdminDeleteOrder(c, ctx)
	})

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/otterEva/lamps/users_service/handlers"
	"github.com/otterEva/lamps/users_service/settings"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "auth_service",
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
		AllowOrigins: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, ctx)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, ctx)
	})

	app.Get("/validattion/:userId/:admin", func(c *fiber.Ctx) error {
		return handlers.CheckForUserHandler(c)
	})

	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

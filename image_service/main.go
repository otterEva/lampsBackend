package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/image_service/handlers"
	"github.com/otterEva/lamps/image_service/settings"
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
	app.Post("/images", func(c *fiber.Ctx) error {
		return handlers.PostImageHandler(c)
	})

	app.Get("/images/:image_url", func(c *fiber.Ctx) error {
		return handlers.GetImageHandler(c)
	})
	// -----------------------------------------------------------------

	if err := app.Listen(settings.Config.APP_URL); err != nil {
		panic(err)
	}
}

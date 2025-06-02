package handlers

import (
	"bytes"
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/otterEva/lamps/goods_service/handlers/helpers"
	"github.com/otterEva/lamps/goods_service/settings"
)

func AdminGoodsPatch(c *fiber.Ctx, ctx context.Context) error {

	val := c.Locals("admin")
	isAdmin, ok := val.(bool)
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	image_url, err := helpers.SendToImageService(c)
	
	log.Debug(image_url)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
	}

	idStr := c.Params("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id := uint(idUint64)

	update := sq.Update("Goods").Where(sq.Eq{"id": id})

	if desc := c.FormValue("description"); desc != "" {
		update = update.Set("description", desc)
	}
	if name := c.FormValue("name"); name != "" {
		update = update.Set("name", name)
	}
	if img := c.FormValue("image_url"); img != "" {
		update = update.Set(image_url, img)
	}
	if costStr := c.FormValue("cost"); costStr != "" {
		if costUint64, err := strconv.ParseUint(costStr, 10, 64); err == nil {
			update = update.Set("cost", uint(costUint64))
		} else {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	if activeStr := c.FormValue("active"); activeStr != "" {
		if activeBool, err := strconv.ParseBool(activeStr); err == nil {
			update = update.Set("active", activeBool)
		} else {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}

	sql, args, err := update.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	dbClient := settings.Clients.DbClient

	_, err = dbClient.Exec(ctx, sql, args...)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

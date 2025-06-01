package handlers

import (
	"bytes"
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/otterEva/lamps/goods_service/handlers/helpers"
	"github.com/otterEva/lamps/goods_service/schemas"
	"github.com/otterEva/lamps/goods_service/settings"
)

func AdminGoodsPost(c *fiber.Ctx, ctx context.Context) error {

	val := c.Locals("admin")
	isAdmin, ok := val.(bool)
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "No file provided")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot open file")
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file")
	}

	image_url, err := helpers.SendToImageService(c)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
	}

	costStr := c.FormValue("cost")
	costUint64, err := strconv.ParseUint(costStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cost",
		})
	}

	activeStr := c.FormValue("active")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid active form",
		})
	}

	good := &schemas.Good{
		Description: c.FormValue("description"),
		Name:        c.FormValue("name"),
		ImageURL:    image_url,
		Cost:        uint(costUint64),
		Active:      active,
	}

	if good.Description == "" || good.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "description, name, and image_url are required",
		})
	}

	sql, args, err := sq.
		Insert("Goods").
		Columns("description", "name", "image_url", "cost", "active").
		Values(good.Description, good.Name, good.ImageURL, good.Cost, good.Active).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Debug(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	dbClient := settings.Clients.DbClient

	_, err = dbClient.Exec(ctx, sql, args...)
	if err != nil {
		log.Debug(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

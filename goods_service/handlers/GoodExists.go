package handlers

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/goods_service/schemas"
	"github.com/otterEva/lamps/goods_service/settings"
)

func CheckIfGoodExists(c *fiber.Ctx, ctx context.Context) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	sql, args, err := sq.
		Select("id", "description", "name", "active", "cost", "image_url").
		From("Goods").
		Where(sq.Eq{"id": id, "active": true}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var g schemas.GoodDB

	dbClient := settings.Clients.DbClient

	err = dbClient.QueryRow(ctx, sql, args...).Scan(&g.ID, &g.Description, &g.Name, &g.Active, &g.Cost, &g.ImageURL)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusOK)
}

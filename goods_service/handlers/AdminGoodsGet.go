package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/otterEva/lamps/goods_service/schemas"
	"github.com/otterEva/lamps/goods_service/settings"
)

func AdminGoodsGet(c *fiber.Ctx, ctx context.Context) error {

	val := c.Locals("admin")
	isAdmin, ok := val.(bool)
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	log.Debug(isAdmin)
		
	sql, args, err := sq.
		Select("id", "description", "name", "image_url", "cost", "active").
		From("Goods").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	dbClient := settings.Clients.DbClient

	rows, err := dbClient.Query(ctx, sql, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var goods []schemas.GoodDB

	for rows.Next() {
		var g schemas.GoodDB
		err := rows.Scan(&g.ID, &g.Description, &g.Name, &g.ImageURL, &g.Cost, &g.Active)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		goods = append(goods, g)
	}

	return c.JSON(goods)
}

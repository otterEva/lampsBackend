package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/goods_service/schemas"
	"github.com/otterEva/lamps/goods_service/settings"
)

func UserGoodsGet(c *fiber.Ctx, ctx context.Context) error {

		sql, args, err := sq.
			Select("id", "description", "name", "active", "cost", "image_url").
			From("Goods").
			Where(sq.Eq{"active": true}).
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
			err := rows.Scan(&g.ID, &g.Description, &g.Name, &g.Active, &g.Cost, &g.ImageURL)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			goods = append(goods, g)
		}

		return c.JSON(goods)
}

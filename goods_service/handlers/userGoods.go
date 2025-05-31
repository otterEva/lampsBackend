package handlers

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otterEva/lamps/goods_service/schemas"
)

func UserGoodsHandlers(route fiber.Router, dbClient *pgxpool.Pool, ctx context.Context) {

	route.Get("/", func(c *fiber.Ctx) error {
		sql, args, err := sq.
			Select("id", "description", "name", "active", "cost", "image_url").
			From("Goods").
			Where(sq.Eq{"active": true}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

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
	})

	// --------------------------------------------------------------------------------------

	route.Get("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}

		sql, args, err := sq.
			Select("id", "description", "name", "active", "cost", "image_url").
			From("Goods").
			Where(sq.Eq{"id": id}).
			Limit(1).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var g schemas.GoodDB
		err = dbClient.QueryRow(ctx, sql, args...).Scan(&g.ID, &g.Description, &g.Name, &g.Active, &g.Cost, &g.ImageURL)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Good not found"})
		}

		return c.JSON(g)
	})
}

package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otterEva/lamps/orders_service/schemas"
)

route.Get("/", func(c *fiber.Ctx) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		sql, args, err := sq.
			Select("order_uuid", "user_id", "good_id", "amount").
			From("Orders").
			OrderBy("order_uuid").
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

		var allOrders []schemas.OrderItem

		for rows.Next() {
			var item schemas.OrderItem
			if err := rows.Scan(&item.OrderUUID, &item.UserID, &item.GoodID, &item.Amount); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			allOrders = append(allOrders, item)
		}

		return c.JSON(allOrders)
	})
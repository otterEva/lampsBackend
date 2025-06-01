package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/orders_service/schemas"
	"github.com/otterEva/lamps/orders_service/settings"
)

func UserGetOrders(c *fiber.Ctx, ctx context.Context) error {
	userIDRaw := c.Locals("userId")
	userID := uint(userIDRaw.(float64))

	sql, args, err := sq.
		Select("order_uuid", "good_id", "amount").
		From("Orders").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("order_uuid").
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

	orders := make(map[string][]schemas.OrderItemOutput)

	for rows.Next() {
		var item schemas.OrderItemOutput
		if err := rows.Scan(&item.OrderUUID, &item.GoodID, &item.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		orders[item.OrderUUID] = append(orders[item.OrderUUID], item)
	}

	return c.JSON(orders)
}

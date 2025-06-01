package handlers

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/orders_service/settings"
)

func AdminDeleteOrder(c *fiber.Ctx, ctx context.Context) error {
	val := c.Locals("admin")
	isAdmin, ok := val.(bool)
	if !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	orderUUID := c.Params("uuid")

	if orderUUID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	sql, args, err := sq.
		Delete("Orders").
		Where(sq.Eq{"order_uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	dbClient := settings.Clients.DbClient

	result, err := dbClient.Exec(ctx, sql, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.RowsAffected() == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.SendStatus(fiber.StatusOK)
}

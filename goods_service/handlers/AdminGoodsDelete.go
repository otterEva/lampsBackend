package handlers

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/goods_service/settings"
)

func AdminGoodDelete(c *fiber.Ctx, ctx context.Context) error {

		val := c.Locals("admin")
		isAdmin, ok := val.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		idStr := c.Params("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		id := uint(idUint64)

		sql, args, err := sq.
			Delete("Goods").
			Where(sq.Eq{"id": id}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}


		dbClient := settings.Clients.DbClient

		cmdTag, err := dbClient.Exec(ctx, sql, args...)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if cmdTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Good not found"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
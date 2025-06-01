package handlers

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/otterEva/lamps/orders_service/schemas"
)

func UserPostOrder(c *fiber.Ctx, ctx context.Context) error {

		var items []schemas.OrderItemInput

		if err := json.Unmarshal(c.Body(), &items); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if len(items) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		userIDRaw := c.Locals("userId")
		userID := uint(userIDRaw.(float64))

		orderUUID := uuid.New().String()

		for _, item := range items {
			if item.Amount == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
			}

			sql, args, err := sq.
				Select("1").
				From("Goods").
				Where(sq.Eq{"id": item.GoodID}).
				Limit(1).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
							return c.SendStatus(fiber.StatusInternalServerError)
			}

			var exists int


			err = dbClient.QueryRow(ctx, sql, args...).Scan(&exists)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Good with ID does not exist",
						"id":    item.GoodID,
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		builder := sq.
			Insert("Orders").
			Columns("order_uuid", "user_id", "good_id", "amount").
			PlaceholderFormat(sq.Dollar)

		for _, item := range items {
			builder = builder.Values(orderUUID, userID, item.GoodID, item.Amount)
		}

		sql, args, err := builder.ToSql()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = dbClient.Exec(ctx, sql, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusCreated)
	}
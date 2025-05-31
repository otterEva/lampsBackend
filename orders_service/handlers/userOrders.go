package handlers

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otterEva/lamps/orders_service/schemas"
)

func UserOrdersHandler(route fiber.Router, dbClient *pgxpool.Pool, ctx context.Context) {

	route.Post("/", func(c *fiber.Ctx) error {
		var items []schemas.OrderItemInput

		if err := json.Unmarshal(c.Body(), &items); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}
		if len(items) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Empty order list",
			})
		}

		userIDRaw := c.Locals("userId")
		userID := uint(userIDRaw.(float64))

		orderUUID := uuid.New().String()

		for _, item := range items {
			if item.Amount == 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Amount must be > 0",
				})
			}

			sql, args, err := sq.
				Select("1").
				From("Goods").
				Where(sq.Eq{"id": item.GoodID}).
				Limit(1).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":    "Order created",
			"order_uuid": orderUUID,
		})
	})

	//---------------------------------------------------------------------------------------------

	route.Get("/", func(c *fiber.Ctx) error {
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
	})
}

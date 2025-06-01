package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/otterEva/lamps/orders_service/schemas"
	"github.com/otterEva/lamps/orders_service/settings"
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

		url := fmt.Sprintf("http://goods_service:8083/%v", item.GoodID)
		resp, err := http.Get(url)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if resp.Status != "200" {
			return c.SendStatus(fiber.StatusBadRequest)
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

	dbClient := settings.Clients.DbClient
	_, err = dbClient.Exec(ctx, sql, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

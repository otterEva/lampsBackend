package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/users_service/utils"
)

func CheckForUserHandler(c *fiber.Ctx) error {
	userId := c.Query("userId")
	admin := c.Query("admin")

	err := utils.GetUserFromDb(context.Background(), userId, admin)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusOK)
}

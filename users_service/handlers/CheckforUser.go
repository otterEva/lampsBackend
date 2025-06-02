package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/users_service/logs"
	"github.com/otterEva/lamps/users_service/utils"
)

func CheckForUserHandler(c *fiber.Ctx) error {

	userIdStr := c.Params("userId")

	id64, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userId, must be a positive integer",
		})
	}
	userId := uint(id64)

	adminStr := c.Params("admin")
	admin, err := strconv.ParseBool(adminStr)

	logs.Logger.Debug("Credentials from c.Params:", "userId", userId, "adminStr", adminStr)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "admin must be ether true ether false",
		})
	}

	err = utils.GetUserFromDb(c, context.Background(), userId, admin)

	if err != nil {
		return err
	}

	logs.Logger.Debug("user exists")
	return c.SendStatus(fiber.StatusOK)
}

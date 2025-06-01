package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/users_service/utils"
)

func CheckForUserHandler(c *fiber.Ctx) error {
	userId := c.Params("userId")
	adminStr := c.Params("admin")

 	admin, err := strconv.ParseBool(adminStr)

    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("admin must be true or false")
    }

	err = utils.GetUserFromDb(context.Background(), userId, admin)

	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

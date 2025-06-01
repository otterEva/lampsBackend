package handlers

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/otterEva/lamps/image_service/utils"
)

func PostImageHandler(c *fiber.Ctx) error {
	val := c.Locals("admin")
	isAdmin, ok := val.(bool)
	if !ok || !isAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	fileName, err := utils.AddFile(fileHeader.Filename, buf.Bytes(), fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(fileName)
}

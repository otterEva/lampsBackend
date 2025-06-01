package handlers

import (
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go"
	"github.com/otterEva/lamps/image_service/settings"
)

func GetImageHandler(c *fiber.Ctx) error {
	imageName := c.Params("image_url")
	if imageName == "" {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	object, err := settings.Clients.MinioClient.GetObject(
		settings.Config.MINIO_BUCKET,
		imageName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	_, err = object.Stat()

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound)
	}

	ext := filepath.Ext(imageName)
	contentType := "application/octet-stream"
	if ext != "" {
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			contentType = mimeType
		}
	}

	c.Set("Content-Type", contentType)

	return c.SendStream(object)
}

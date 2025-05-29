package handlers

import (
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go"
	"github.com/otterEva/lamps/app/settings"
)

func GetImageHandler(c *fiber.Ctx) error {
	imageName := c.Params("image_url")
	if imageName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "image_url path parameter is required")
	}

	object, err := settings.Clients.MinioClient.GetObject(
		settings.Config.MINIO_BUCKET,
		imageName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get object from MinIO")
	}

	_, err = object.Stat()
	
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "image not found")
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

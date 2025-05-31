package handlers

import (
	"bytes"
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// -------------------------------------------------------------------

func PostImageHandler(c *fiber.Ctx, originalFilename string, data []byte, contentType string) (string, error) {
	if contentType != "image/jpeg" && contentType != "image/png" {
		return "", fiber.NewError(fiber.StatusBadRequest)
	}

	ext := filepath.Ext(originalFilename)
	if ext == "" {
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		} else {
			return "", fiber.NewError(fiber.StatusBadRequest)
		}
	}

	objectName := uuid.New().String() + ext
	reader := bytes.NewReader(data)

	_, err := settings.Clients.MinioClient.PutObject(
		settings.Config.MINIO_BUCKET,
		objectName,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError)
	}

	return objectName, nil
}
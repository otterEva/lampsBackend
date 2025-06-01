package utils

import (
	"bytes"
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
	"github.com/otterEva/lamps/image_service/settings"
)

func AddFile(originalFilename string, data []byte, contentType string) (string, error) {
	if contentType != "image/jpeg" && contentType != "image/png" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Only JPEG and PNG allowed")
	}

	ext := filepath.Ext(originalFilename)
	if ext == "" {
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		} else {
			return "", fiber.NewError(fiber.StatusBadRequest, "Invalid file type")
		}
	}

	objectName := uuid.New().String() + ext
	reader := bytes.NewReader(data)

	log.Debug(settings.Config.MINIO_BUCKET, objectName, contentType)

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
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file")
	}

	return objectName, nil
}

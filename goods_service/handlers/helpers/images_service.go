package helpers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SendToImageService(c *fiber.Ctx) (string, error) {

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return "", c.SendStatus(fiber.StatusBadRequest)
		}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
		}

	defer srcFile.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
		}

	if _, err := io.Copy(part, srcFile); err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := writer.Close(); err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}

	imageServiceURL := "http://image_service:8003/images"

	req, err := http.NewRequest("POST", imageServiceURL, &buf)
	if err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", c.SendStatus(fiber.StatusBadGateway)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", c.SendStatus(fiber.StatusBadGateway)
	}

	savedNameBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}

	return string(savedNameBytes), nil
}

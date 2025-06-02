package helpers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

    contentType := fileHeader.Header.Get("Content-Type")
    if contentType == "" {

        switch {
        case fileHeader.Filename[len(fileHeader.Filename)-4:] == ".jpg",
            fileHeader.Filename[len(fileHeader.Filename)-5:] == ".jpeg":
            contentType = "image/jpeg"
        case fileHeader.Filename[len(fileHeader.Filename)-4:] == ".png":
            contentType = "image/png"
        default:
            contentType = "application/octet-stream"
        }
    }

    partHeader := make(textproto.MIMEHeader)
    partHeader.Set("Content-Disposition",
        fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileHeader.Filename))
    partHeader.Set("Content-Type", contentType)

    part, err := writer.CreatePart(partHeader)
    if err != nil {
        return "", c.SendStatus(fiber.StatusInternalServerError)
    }

    if _, err := io.Copy(part, srcFile); err != nil {
        return "", c.SendStatus(fiber.StatusInternalServerError)
    }
    if err := writer.Close(); err != nil {
        return "", c.SendStatus(fiber.StatusInternalServerError)
    }

    imageServiceURL := "http://image_service:8080/images"
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
    log.Debug(resp)

    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        return "", c.SendStatus(fiber.StatusBadGateway)
    }

    savedNameBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", c.SendStatus(fiber.StatusInternalServerError)
    }
    log.Debug(savedNameBytes)

    return string(savedNameBytes), nil
}

package api

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type WebhookData struct {
	PushData struct {
		Tag       string `json:"tag"`
		MediaType string `json:"media_type"`
	} `json:"push_data"`
}

func WebhookHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != http.MethodPost {
			log.Println("Method Not Allowed")
			return c.Status(http.StatusMethodNotAllowed).SendString("Method Not Allowed")
		}

		var webhookData WebhookData
		if err := c.BodyParser(&webhookData); err != nil {
			log.Println("Failed to parse webhook data:", err)
			return c.Status(http.StatusBadRequest).SendString("Bad Request")
		}

		log.Println("Webhook received:", webhookData)

		return c.SendString("Success")
	}
}

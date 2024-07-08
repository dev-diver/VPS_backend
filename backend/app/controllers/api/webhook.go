package api

import (
	"bytes"
	"fmt"
	"io"
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
		var webhookData WebhookData
		if err := c.BodyParser(&webhookData); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		log.Printf("Webhook received: %+v", webhookData)

		//client stop
		if err := dockerRequest("POST", "/containers/client/stop", nil); err != nil {
			log.Printf("Failed to stop container: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		//client rm -f
		if err := dockerRequest("DELETE", "/containers/client", nil); err != nil {
			log.Printf("Failed to remove container: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		//compose pull client
		imageName := "client:latest"
		pullImageUrl := fmt.Sprintf("/images/create?fromImage=%s", imageName)
		if err := dockerRequest("POST", pullImageUrl, nil); err != nil {
			log.Printf("Failed to pull images: %v", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		//compose run -d client
		runContainerData := []byte(fmt.Sprintf(`{"Image": "%s"}`, imageName))
		if err := dockerRequest("POST", "/containers/create", runContainerData); err != nil {
			log.Printf("Failed to create container: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		// //docker restart server
		// if err := dockerRequest("POST", "/containers/server/restart", nil); err != nil {
		// 	log.Printf("Failed to restart container: %v", err)
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		// }

		return c.JSON(fiber.Map{"status": "Success"})
	}
}

func dockerRequest(method, command string, jsonData []byte) error {
	client := &http.Client{}
	url := "http://localhost:2375" // Docker REST API endpoint

	var req *http.Request
	var err error
	if jsonData != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create container: %s", body)
	}
	return nil
}

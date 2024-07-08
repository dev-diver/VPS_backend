package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type WebhookData struct {
	PushData struct {
		Tag       string `json:"tag"`
		MediaType string `json:"media_type"`
	} `json:"push_data"`
	Repository struct {
		RepoName string `json:"repo_name"`
	} `json:"repository"`
}

func WebhookHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var webhookData WebhookData
		if err := c.BodyParser(&webhookData); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		log.Printf("Webhook received: %+v", webhookData)

		if imageName := webhookData.Repository.RepoName; imageName == "devdiver/vacation_promotion_client" {
			if err := clientRestart(imageName); err != nil {
				return err
			}
		} else if imageName == "devdiver/vacation_promotion_server" {
			if err := serverRestart(imageName); err != nil {
				return err
			}
		} else {
			log.Printf("Unknown repository: %s", webhookData.Repository.RepoName)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid repository"})
		}

		return c.JSON(fiber.Map{"status": "Success"})
	}
}

func clientRestart(imageName string) error {
	//client stop
	if err := dockerRequest("POST", "/containers/client/stop", nil); err != nil {
		log.Printf("Failed to stop client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	//client rm -f
	if err := dockerRequest("DELETE", "/containers/client", nil); err != nil {
		log.Printf("Failed to remove client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	//compose pull client
	if err := imagePull(imageName); err != nil {
		return err
	}

	//compose run -d client
	runContainerData := []byte(`{
		"Image": "devdiver/vacation_promotion_client:latest",
		"HostConfig": {
			"Binds": ["/front_app:/dist"],
			"Command": ["sh", "-c", "rm -rf /dist/* && mv /app/front_web/dist/* /dist && bin/true"]
		}
	}`)
	if err := dockerRequest("POST", "/containers/create", runContainerData); err != nil {
		log.Printf("Failed to create container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := dockerRequest("POST", "/containers/vacation_promotion_client/start", nil); err != nil {
		log.Printf("Failed to start client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func serverRestart(imageName string) error {

	//image pull server
	if err := imagePull(imageName); err != nil {
		return err
	}

	//docker restart server
	if err := dockerRequest("POST", "/containers/server/restart", nil); err != nil {
		log.Printf("Failed to restart container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	return nil
}

func imagePull(imageName string) error {
	pullImageUrl := fmt.Sprintf("/images/create?fromImage=%s:latest", imageName)
	if err := dockerRequest("POST", pullImageUrl, nil); err != nil {
		log.Printf("Failed to pull images: %v", nil)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	return nil
}

func dockerRequest(method, command string, jsonData []byte) error {
	client := &http.Client{}
	hostIP := os.Getenv("HOST_IP") // 환경 변수에서 호스트 IP 주소 가져오기

	url := "http://" + hostIP + ":2375" + command
	fmt.Printf("request to %s", url)

	var req *http.Request
	var err error
	if jsonData != nil {
		req, err = http.NewRequest(method, url+command, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url+command, nil)
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

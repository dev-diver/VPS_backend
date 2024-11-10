package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

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

func UpdateHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var webhookData WebhookData
		if err := c.BodyParser(&webhookData); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		log.Printf("Webhook received: %+v", webhookData)

		if webhookData.PushData.Tag != "latest" {
			log.Printf("Tag is not latest: %s", webhookData.PushData.Tag)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag"})
		}

		if imageName := webhookData.Repository.RepoName; imageName == "devdiver/vacation_promotion_client" {
			log.Println("Sending update request for client...")
			if err := sendUpdateRequest("client"); err != nil {
				log.Printf("Error updating client: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		} else if imageName == "devdiver/vacation_promotion_server" {
			log.Println("Sending update request for server...")
			if err := sendUpdateRequest("server"); err != nil {
				log.Printf("Error updating server: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			log.Printf("Unknown repository: %s", webhookData.Repository.RepoName)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid repository"})
		}

		return c.JSON(fiber.Map{"status": "Success"})
	}
}

func sendUpdateRequest(serviceName string) error {
	url := "http://update-server:5000/update"

	postData := map[string]string{
		"service_name": serviceName,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return errors.New("error marshalling JSON: " + err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("error sending update request: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error reading response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update " + serviceName + ": " + string(body))
	}

	log.Printf("%s service update triggered successfully", serviceName)
	log.Printf("message: %s", body)
	return nil
}

func HaveUpdateHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceName := c.Query("service")
		if serviceName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var imageName string
		if serviceName == "server" {
			imageName = "devdiver/vacation_promotion_server"
		} else if serviceName == "client" {
			imageName = "devdiver/vacation_promotion_client"
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service"})
		}

		remoteCreated, err := getRemoteImageCreatedTime(imageName)
		if err != nil {
			return err
		}
		log.Printf("Latest digest for %s: %s", imageName, remoteCreated)

		localCreated, err := getLocalImageCreatedTime(imageName)
		if err != nil {
			return err
		}
		log.Printf("Current digest for %s: %s", serviceName, localCreated)

		if remoteCreated.After(localCreated) {
			return c.JSON(fiber.Map{
				"update":  true,
				"latest":  remoteCreated,
				"current": localCreated,
			})
		} else {
			return c.JSON(fiber.Map{
				"update":  false,
				"latest":  remoteCreated,
				"current": localCreated,
			})
		}
	}
}

func getRemoteImageCreatedTime(imageName string) (time.Time, error) {
	cmd := exec.Command("regctl", "image", "inspect", imageName, "--format", "{{.Created}}")
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}
	const layout = "2006-01-02 15:04:05.999999999 -0700 MST"

	// 시간 문자열을 파싱
	createdTime, err := time.Parse(layout, strings.TrimSpace(string(output)))
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}

	return createdTime, nil
}

func getLocalImageCreatedTime(imageName string) (time.Time, error) {
	cmd := exec.Command("docker", "inspect", "--format={{.Created}}", imageName)
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}
	log.Printf("Inspect output: %s", output)

	// 출력된 시간 문자열 파싱
	timeStr := strings.Trim(string(output), "'\n") // 작은 따옴표와 줄바꿈 문자 제거
	createdTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time: %v", err)
	}

	return createdTime, nil
}

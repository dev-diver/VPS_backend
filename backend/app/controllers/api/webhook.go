package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

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

		if webhookData.PushData.Tag != "latest" {
			log.Printf("Tag is not latest: %s", webhookData.PushData.Tag)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag"})
		}

		if imageName := webhookData.Repository.RepoName; imageName == "devdiver/vacation_promotion_client" {
			// if err := clientRestart(imageName); err != nil {
			// 	return err
			// }
			if err := clientRestartWithSocket(); err != nil {
				return err
			}
		} else if imageName == "devdiver/vacation_promotion_server" {
			// if err := serverRestart(imageName); err != nil {
			// 	return err
			// }
			if err := serverRestartWithSocket(); err != nil {
				return err
			}
		} else {
			log.Printf("Unknown repository: %s", webhookData.Repository.RepoName)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid repository"})
		}

		return c.JSON(fiber.Map{"status": "Success"})
	}
}

func execCommand(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func clientRestartWithSocket() error {

	containerName := "vacation_promotion_client"
	newName := containerName + "_old"
	log.Printf("rename client")
	if err := exec.Command("docker", "rename", containerName, newName).Run(); err != nil {
		log.Printf("Failed to stop client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker compose pull client")
	if err := execCommand("docker", "compose", "pull", "client"); err != nil {
		log.Printf("Failed to pull client: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker compose up -d client ")
	if err := execCommand("docker", "compose", "up", "client", "-d"); err != nil {
		log.Printf("Failed to start client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker rm -f client ")
	if err := execCommand("docker", "rm", "-f", newName); err != nil {
		log.Printf("Failed to remove client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func serverRestartWithSocket() error {

	containerName := "vacation_promotion_server"
	newName := containerName + "_old"
	log.Printf("rename server")
	if err := exec.Command("docker", "rename", containerName, newName).Run(); err != nil {
		log.Printf("Failed to stop server container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker compose pull server")
	if err := execCommand("docker", "compose", "pull", "server"); err != nil {
		log.Printf("Failed to pull server: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker compose up -d server")
	if err := execCommand("docker", "compose", "up", "server", "-d"); err != nil {
		log.Printf("Failed to start server container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	log.Printf("docker rm -f server ")
	if err := execCommand("docker", "rm", "-f", newName); err != nil {
		log.Printf("Failed to remove server container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func clientRestart(imageName string) error {

	client_container_name := "vacation_promotion_client"

	//client stop
	if err := dockerRequest("POST", fmt.Sprintf("/containers/%s/stop", client_container_name), nil); err != nil {
		log.Printf("Failed to stop client container: %v", err)
	}

	//client rm -f
	if err := dockerRequest("DELETE", fmt.Sprintf("/containers/%s", client_container_name), nil); err != nil {
		log.Printf("Failed to remove client container: %v", err)
	}

	//compose pull client
	if err := imagePull(imageName); err != nil {
		log.Printf("Failed to pull client %v", err)
		return err
	}

	//compose run -d client
	runContainerData := []byte(`{
		"Image": "devdiver/vacation_promotion_client:latest",
		"HostConfig": {
			"Binds": ["vps_central_front_app:/dist"],
			"Command": ["sh", "-c", "rm -rf /dist/* && mv /app/front_web/dist/* /dist && bin/true"]
		}
	}`)
	if err := dockerRequest("POST", fmt.Sprintf("/containers/create?name=%s", client_container_name), runContainerData); err != nil {
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
		log.Printf("Failed to pull server %v", err)
		return err
	}

	server_container_name := "vacation_promotion_server"

	// log.Println("Creating and starting server container...")
	// runContainerData := []byte(`{
	// 	"Image": "devdiver/vacation_promotion_server:latest",
	// 	"Env": ["HOST_IP=${HOST_IP}"],
	// 	"HostConfig": {
	// 		"Binds": [
	// 			"config:/app/backend/config",
	// 			"database:/app/backend/database",
	// 			"front_app:/app/dist"
	// 		],
	// 		"PortBindings": {
	// 			"3000/tcp": [{"HostPort": "3000"}]
	// 		}
	// 	}
	// }`)
	// if err := dockerRequest("POST", fmt.Sprintf("/containers/create?name=%s", server_container_name), runContainerData); err != nil {
	// 	log.Printf("Failed to create server container: %v", err)
	// 	return err
	// }

	// if err := dockerRequest("POST", fmt.Sprintf("/containers/%s/start", server_container_name), nil); err != nil {
	// 	log.Printf("Failed to start server container: %v", err)
	// 	return err
	// }
	// log.Println("Server container created and started successfully.")

	//docker restart server
	if err := dockerRequest("POST", fmt.Sprintf("/containers/%s/restart?t=5", server_container_name), nil); err != nil {
		log.Printf("Failed to restart container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	log.Println("Server container restarted successfully.")
	return nil
}

func imagePull(imageName string) error {
	pullImageUrl := fmt.Sprintf("/images/create?fromImage=%s:latest", url.QueryEscape(imageName))
	if err := dockerRequest("POST", pullImageUrl, nil); err != nil {
		log.Printf("Failed to pull images: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	return nil
}

func dockerRequest(method, command string, jsonData []byte) error {
	client := &http.Client{}
	hostIP := os.Getenv("HOST_IP") // 환경 변수에서 호스트 IP 주소 가져오기

	url := "http://" + hostIP + ":2375" + command
	log.Printf("%s request to %s", method, url)

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

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker API error: %s", body)
	}

	log.Printf("Request to %s completed with status %d", url, resp.StatusCode)
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

		latestDigest, err := getLatestDockerDigest(imageName)
		if err != nil {
			return err
		}
		log.Printf("Latest digest for %s: %s", imageName, latestDigest)

		currentDigest, err := getCurrentImageDigest(serviceName)
		if err != nil {
			return err
		}
		log.Printf("Current digest for %s: %s", serviceName, currentDigest)

		if currentDigest != latestDigest {
			return c.JSON(fiber.Map{"update": true})
		} else {
			return c.JSON(fiber.Map{"update": false})
		}
	}
}

type Descriptor struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
	Platform  struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
	} `json:"platform"`
}
type Manifest struct {
	Ref         string     `json:"Ref"`
	Descriptor  Descriptor `json:"Descriptor"`
	Raw         string     `json:"Raw"`
	OCIManifest struct {
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
		Config        struct {
			MediaType string `json:"mediaType"`
			Digest    string `json:"digest"`
			Size      int    `json:"size"`
		} `json:"config"`
		Layers []struct {
			MediaType   string `json:"mediaType"`
			Digest      string `json:"digest"`
			Size        int    `json:"size"`
			Annotations struct {
				PredicateType string `json:"in-toto.io/predicate-type"`
			} `json:"annotations,omitempty"`
		} `json:"layers"`
	} `json:"OCIManifest"`
}

func getLatestDockerDigest(image string) (string, error) {
	cmd := exec.Command("docker", "manifest", "inspect", image+":latest", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to inspect manifest: %v", err)
	}

	log.Printf("Manifest: %s", out.String())

	var manifests []Manifest
	if err := json.Unmarshal(out.Bytes(), &manifests); err != nil {
		return "", fmt.Errorf("failed to unmarshal manifest: %v", err)
	}

	if len(manifests) > 0 {
		return manifests[0].Descriptor.Digest, nil
	}
	return "", fmt.Errorf("no manifests found")
}

func getCurrentImageDigest(containerName string) (string, error) {
	cmd := exec.Command("docker", "inspect", "--format='{{index .RepoDigests 0}}'", containerName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(output), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected output: %s", output)
	}
	return strings.TrimSpace(parts[1]), nil
}

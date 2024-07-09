package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		// 기본적으로 sh 사용
		shell = "/bin/bash"
	}
	return shell
}

func execCommand(command string, args ...string) error {
	shell := getShell()
	fullCommand := append([]string{"-c", command}, args...)
	cmd := exec.Command(shell, fullCommand...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Running command: %v", cmd)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func clientRestartWithSocket() error {

	containerName := "vacation_promotion_client"
	newName := containerName + "_old"
	if err := execCommand("docker", "rename", containerName, newName); err != nil {
		log.Printf("Failed to stop client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := execCommand("docker", "compose", "pull", "client"); err != nil {
		log.Printf("Failed to pull client: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := execCommand("docker", "compose", "up", "client", "-d"); err != nil {
		log.Printf("Failed to start client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := execCommand("docker", "rm", "-f", newName); err != nil {
		log.Printf("Failed to remove client container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func serverRestartWithSocket() error {

	// containerName := "vacation_promotion_server"
	// newName := containerName + "_old"
	// if err := execCommand("docker", "rename", containerName, newName); err != nil {
	// 	log.Printf("Failed to stop server container: %v", err)
	// 	return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	// }

	if err := execCommand("docker", "compose", "pull", "server"); err != nil {
		log.Printf("Failed to pull server: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := execCommand("docker", "compose", "up", "server", "-d"); err != nil {
		log.Printf("Failed to start server container: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// if err := execCommand("docker", "rm", "-f", newName); err != nil {
	// 	log.Printf("Failed to remove server container: %v", err)
	// 	return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	// }

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

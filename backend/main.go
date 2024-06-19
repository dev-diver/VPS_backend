package main

import (
	"log"

	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config, err := database.LoadConfig("database/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New()

	database.InitDB(config)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}

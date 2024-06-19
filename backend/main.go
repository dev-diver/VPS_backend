package main

import (
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	database.InitDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}

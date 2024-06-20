package api

import (
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement login logic here
		return c.SendStatus(fiber.StatusOK)
	}
}

func LogoutHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement logout logic here
		return c.SendStatus(fiber.StatusOK)
	}
}

package api

import "github.com/gofiber/fiber/v2"

func LoginHandler(c *fiber.Ctx) error {
	// Implement login logic here
	return c.SendStatus(fiber.StatusOK)
}

func LogoutHandler(c *fiber.Ctx) error {
	// Implement logout logic here
	return c.SendStatus(fiber.StatusOK)
}

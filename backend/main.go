package main

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	database.InitDB()
	database.DBConn.AutoMigrate(
		&models.Company{},
		&models.Member{},
		&models.MemberAdmin{},
		&models.NotificationMember{},
		&models.Group{},
		&models.GivenVacation{},
		&models.ApplyVacation{},
		&models.VacationPlan{},
		&models.Notification{},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}

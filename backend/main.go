package main

import (
	"fmt"
	"log"

	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"cywell.com/vacation-promotion/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(
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
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	fmt.Println("Database migrated successfully")

	api := app.Group("/api")
	routes.RegisterAPI(api, db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}

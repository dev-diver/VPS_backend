package api

import (
	"time"

	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func ApplyVacationHandler(c *fiber.Ctx, db *database.Database) error {
	vacation := new(models.ApplyVacation)
	if err := c.BodyParser(vacation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	vacation.VacationPlan.ApproveDate = time.Now() // Example: setting approve date
	if err := db.Create(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(vacation)
}

func CancelVacationHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("vacationID")
	if err := db.Delete(&models.ApplyVacation{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func ApproveVacationHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("vacationID")
	var vacation models.ApplyVacation
	if err := db.First(&vacation, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	vacation.VacationPlan.ApproveDate = time.Now() // Example: setting approve date
	if err := db.Save(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacation)
}

func RejectVacationHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("vacationID")
	var vacation models.ApplyVacation
	if err := db.First(&vacation, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	vacation.VacationPlan.VacationProcessStateID = 3 // Assuming 3 is the ID for the rejected state
	if err := db.Save(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacation)
}

func GetMemberVacationsHandler(c *fiber.Ctx, db *database.Database) error {
	memberID := c.Params("memberID")
	year := c.Params("year")
	month := c.Params("month")
	var vacations []models.GivenVacation

	query := db.Where("member_id = ? AND year = ?", memberID, year)
	if month != "" {
		query = query.Where("MONTH(generate_date) = ?", month)
	}

	if err := query.Find(&vacations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacations)
}

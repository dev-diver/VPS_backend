package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func CreateCompanyHandler(c *fiber.Ctx, db *database.Database) error {
	company := new(models.Company)
	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(company)
}

func DeleteCompanyHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("companyID")
	if err := db.Delete(&models.Company{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func GetCompanyHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("companyID")
	var company models.Company
	if err := db.Preload("Members").Preload("Groups").First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(company)
}

func UpdateCompanyHandler(c *fiber.Ctx, db *database.Database) error {

	id := c.Params("companyID")
	var company models.Company
	if err := db.First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Save(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(company)
}

func GetVacationsByYearMonthHandler(c *fiber.Ctx, db *database.Database) error {

	companyID := c.Params("companyID")
	year := c.Params("year")
	month := c.Params("month")
	var vacations []models.GivenVacation

	query := db.Where("company_id = ? AND year = ?", companyID, year)
	if month != "" {
		query = query.Where("MONTH(generate_date) = ?", month)
	}

	if err := query.Find(&vacations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacations)
}

func CreateVacationHandler(c *fiber.Ctx, db *database.Database) error {

	vacation := new(models.GivenVacation)
	if err := c.BodyParser(vacation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(vacation)
}

func UpdateVacationHandler(c *fiber.Ctx, db *database.Database) error {

	id := c.Params("vacationID")
	var vacation models.GivenVacation
	if err := db.First(&vacation, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.BodyParser(&vacation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Save(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacation)
}

func DeleteVacationHandler(c *fiber.Ctx, db *database.Database) error {

	id := c.Params("vacationID")
	if err := db.Delete(&models.GivenVacation{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func PromoteVacationHandler(c *fiber.Ctx, db *database.Database) error {

	id := c.Params("vacationID")
	var vacation models.GivenVacation
	if err := db.First(&vacation, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	vacation.VacationPromotionStateID = 2 // Assuming 2 is the ID for the promotion state
	if err := db.Save(&vacation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(vacation)
}

func GetPromotionsHandler(c *fiber.Ctx, db *database.Database) error {

	companyID := c.Params("companyID")
	var promotions []models.GivenVacation
	if err := db.Where("company_id = ? AND vacation_promotion_state_id = ?", companyID, 2).Find(&promotions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(promotions)
}

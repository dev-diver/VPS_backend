package api

import (
	"errors"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCompany(tx *gorm.DB, company models.Company) (*models.Company, *models.Organize, error) {

	if tx.Error != nil {
		return nil, nil, errors.New("could not begin transaction")
	}

	// 회사 생성
	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("could not create company")
	}

	// 조직 생성
	organize := models.Organize{
		CompanyID: company.ID,
		Name:      company.Name,
		ParentID:  nil,
	}

	if err := tx.Create(&organize).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("could not create organize")
	}

	return &company, &organize, nil
}

func CreateCompanyHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		company := new(models.Company)
		if err := c.BodyParser(company); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		tx := db.DB.Begin()
		company, _, err := CreateCompany(tx, *company)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(company)
	}
}

func GetCompanyHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("companyID")
		var company models.Company
		if err := db.DB.Preload("VacationGenerateType").First(&company, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// company를 companyResponse형으로 변환
		companyResponse := dto.CompanyResponse{
			ID:                          company.ID,
			Name:                        company.Name,
			AccountingDay:               company.AccountingDay,
			VacationGenerateTypeName:    company.VacationGenerateType.TypeName,
			VacationGenerateDescription: company.VacationGenerateType.Description,
		}

		return c.JSON(companyResponse)
	}
}

func UpdateCompanyHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("companyID")
		var company models.Company
		if err := db.DB.Preload("VacationGenerateType").First(&company, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if err := c.BodyParser(&company); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Save(&company).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(company)
	}
}

func DeleteCompanyHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("companyID")
		if err := db.DB.Delete(&models.Company{}, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetCompanyMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyID")
		var members []*models.Member
		if err := db.DB.Where("company_id = ?", companyID).Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(dto.MapMembersToDTO(members))
	}
}

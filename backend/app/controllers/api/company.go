package api

import (
	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func CreateCompanyHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		company := new(models.Company)
		if err := c.BodyParser(company); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Create(&company).Error; err != nil {
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
		var members []models.Member
		if err := db.DB.Where("company_id = ?", companyID).Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var membersResponse []dto.MemberResponse
		for _, member := range members {
			memberResponse := dto.MemberResponse{
				ID:       member.ID,
				Name:     member.Name,
				Email:    member.Email,
				HireDate: member.HireDate,
				IsActive: member.IsActive,
			}
			membersResponse = append(membersResponse, memberResponse)
		}

		return c.JSON(membersResponse)
	}
}

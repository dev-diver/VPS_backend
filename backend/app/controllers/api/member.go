package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func GetMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyID")
		var members []models.Member
		if err := db.DB.Where("company_id = ?", companyID).Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(members)
	}
}

func SearchMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyword := c.Query("keyword")
		companyID := c.Params("companyID")
		var members []models.Member
		if err := db.DB.Where("company_id = ? AND (name LIKE ? OR email LIKE ?)", companyID, "%"+keyword+"%", "%"+keyword+"%").Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(members)
	}
}

func CreateMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		members := new([]models.Member)
		if err := c.BodyParser(members); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Create(members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(members)
	}
}

func GetMemberProfileHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("memberID")
		var member models.Member
		if err := db.DB.First(&member, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(member)
	}
}

func DeactivateMemberHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("memberID")
		var member models.Member
		if err := db.DB.First(&member, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		member.IsActive = false
		if err := db.DB.Save(&member).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(member)
	}
}

func DeleteMemberHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("memberID")
		if err := db.DB.Delete(&models.Member{}, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

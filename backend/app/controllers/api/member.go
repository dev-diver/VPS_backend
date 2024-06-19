package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func GetMembersHandler(c *fiber.Ctx, db *database.Database) error {

	companyID := c.Params("companyID")
	var members []models.Member
	if err := db.Where("company_id = ?", companyID).Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(members)
}

func SearchMembersHandler(c *fiber.Ctx, db *database.Database) error {
	keyword := c.Query("keyword")
	companyID := c.Params("companyID")
	var members []models.Member
	if err := db.Where("company_id = ? AND (name LIKE ? OR email LIKE ?)", companyID, "%"+keyword+"%", "%"+keyword+"%").Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(members)
}

func CreateMembersHandler(c *fiber.Ctx, db *database.Database) error {
	members := new([]models.Member)
	if err := c.BodyParser(members); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(members)
}

func GetMemberProfileHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("memberID")
	var member models.Member
	if err := db.First(&member, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(member)
}

func DeactivateMemberHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("memberID")
	var member models.Member
	if err := db.First(&member, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	member.IsActive = false
	if err := db.Save(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(member)
}

func DeleteMemberHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("memberID")
	if err := db.Delete(&models.Member{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

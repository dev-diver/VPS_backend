package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func GetGroupsHandler(c *fiber.Ctx, db *database.Database) error {
	companyID := c.Params("companyID")
	var groups []models.Group
	if err := db.Preload("Members").Where("company_id = ?", companyID).Find(&groups).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

func CreateGroupHandler(c *fiber.Ctx, db *database.Database) error {
	group := new(models.Group)
	if err := c.BodyParser(group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

func UpdateGroupHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("groupID")
	var group models.Group
	if err := db.First(&group, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Save(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(group)
}

func DeleteGroupHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("groupID")
	if err := db.Delete(&models.Group{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func UpdateGroupMembersHandler(c *fiber.Ctx, db *database.Database) error {
	id := c.Params("groupID")
	var group models.Group
	if err := db.First(&group, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var members []models.Member
	if err := c.BodyParser(&members); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.Model(&group).Association("Members").Replace(&members); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(group)
}

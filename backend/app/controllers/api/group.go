package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func GetGroupsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyID")
		var groups []models.Group
		if err := db.DB.Preload("Members").Where("company_id = ?", companyID).Find(&groups).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(groups)
	}
}

func CreateGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		group := new(models.Group)
		if err := c.BodyParser(group); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Create(&group).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(group)
	}
}

func UpdateGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("groupID")
		var group models.Group
		if err := db.DB.First(&group, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if err := c.BodyParser(&group); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Save(&group).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(group)
	}
}

func DeleteGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("groupID")
		if err := db.DB.Delete(&models.Group{}, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateGroupMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("groupID")
		var group models.Group
		if err := db.DB.First(&group, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var members []models.Member
		if err := c.BodyParser(&members); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := db.DB.Model(&group).Association("Members").Replace(&members); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(group)
	}
}

func GetGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupID := c.Params("groupID")
		var group models.Group
		if err := db.DB.Preload("Members").First(&group, groupID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(group)
	}
}

func AddGroupMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupID := c.Params("groupID")
		var group models.Group
		if err := db.DB.First(&group, groupID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var memberIDs []uint
		if err := c.BodyParser(&memberIDs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var members []models.Member
		if err := db.DB.Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := db.DB.Model(&group).Association("Members").Append(&members); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(group)
	}
}

func DeleteGroupMemberHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupID := c.Params("groupID")
		memberID := c.Params("memberID")
		var group models.Group
		if err := db.DB.First(&group, groupID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var member models.Member
		if err := db.DB.First(&member, memberID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := db.DB.Model(&group).Association("Members").Delete(&member); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

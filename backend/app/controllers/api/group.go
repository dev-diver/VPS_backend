package api

import (
	"strconv"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		companyID, err := strconv.ParseUint(c.Params("companyID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
		}

		var createGroupDTO dto.CreateGroupDTO

		// 요청 바디 파싱
		if err := c.BodyParser(&createGroupDTO); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 유효성 검사
		validate := validator.New()
		if err := validate.Struct(&createGroupDTO); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		priority := 1
		if createGroupDTO.Priority != nil {
			priority = *createGroupDTO.Priority
		}

		group := models.Group{
			CompanyID: uint(companyID),
			Name:      createGroupDTO.Name,
			Color:     createGroupDTO.Color,
			Priority:  priority,
		}

		if err := db.DB.Create(&group).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		groupDTO := dto.GroupDTO{
			ID:        group.ID,
			CompanyID: group.CompanyID,
			Name:      group.Name,
			Color:     group.Color,
			Priority:  group.Priority,
			Members:   group.Members,
		}

		return c.Status(fiber.StatusCreated).JSON(groupDTO)
	}
}

func GetGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupID := c.Params("groupID")
		var group dto.GroupDTO
		if err := db.DB.Preload("Members").First(&group, groupID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(group)
	}
}

func GetGroupsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyID")
		var groups []dto.GroupDTO
		if err := db.DB.Preload("Members").Where("company_id = ?", companyID).Find(&groups).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(groups)
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

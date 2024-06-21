package api

import (
	"strconv"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		companyID, err := strconv.ParseUint(c.Params("companyID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
		}

		var createGroupDTO dto.CreateGroupRequest

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

		groupDTO := dto.GroupResponse{
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
		var group models.Group
		if err := db.DB.Preload("Members").First(&group, "id = ?", groupID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// GroupDTO로 변환합니다.
		groupDTO := dto.GroupResponse{
			ID:        group.ID,
			CompanyID: group.CompanyID,
			Name:      group.Name,
			Color:     group.Color,
			Priority:  group.Priority,
			Members:   group.Members,
		}
		return c.JSON(groupDTO)
	}
}

func GetGroupsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, err := strconv.ParseUint(c.Params("companyID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
		}
		var groups []models.Group
		if err := db.DB.Where("company_id = ?", companyID).Find(&groups).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Groups를 GroupDTO로 변환합니다.
		var groupDTOs []dto.GroupResponse
		for _, group := range groups {
			groupDTO := dto.GroupResponse{
				ID:        group.ID,
				CompanyID: group.CompanyID,
				Name:      group.Name,
				Color:     group.Color,
				Priority:  group.Priority,
				Members:   group.Members,
			}
			groupDTOs = append(groupDTOs, groupDTO)
		}

		return c.JSON(groupDTOs)
	}
}

func UpdateGroupHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("groupID")

		// 기존 그룹을 로드
		var group models.Group
		if err := db.DB.First(&group, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 요청 바디 파싱
		var updateGroupDTO dto.CreateGroupRequest
		if err := c.BodyParser(&updateGroupDTO); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 유효성 검사
		validate := validator.New()
		if err := validate.Struct(&updateGroupDTO); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		// DTO의 데이터를 기존 그룹에 덮어쓰기.
		group.Name = updateGroupDTO.Name
		group.Color = updateGroupDTO.Color
		if updateGroupDTO.Priority != nil {
			group.Priority = *updateGroupDTO.Priority
		}

		// 변경된 내용을 저장
		if err := db.DB.Save(&group).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// GroupDTO로 변환
		groupDTO := dto.GroupResponse{
			ID:        group.ID,
			CompanyID: group.CompanyID,
			Name:      group.Name,
			Color:     group.Color,
			Priority:  group.Priority,
			Members:   group.Members,
		}

		return c.JSON(groupDTO)
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

		var memberIDs []uint
		if err := c.BodyParser(&memberIDs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 멤버 ID 배열을 통해 멤버를 조회
		var members []models.Member
		if err := db.DB.Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 그룹과 멤버 간의 관계를 업데이트
		if err := db.DB.Model(&group).Association("Members").Replace(&members); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 멤버 정보를 MemberDTO로 변환
		var memberDTOs []dto.MemberResponse
		for _, member := range members {
			memberDTO := dto.MemberResponse{
				ID:       member.ID,
				Name:     member.Name,
				Email:    member.Email,
				HireDate: member.HireDate,
				IsActive: member.IsActive,
			}
			memberDTOs = append(memberDTOs, memberDTO)
		}

		return c.JSON(memberDTOs)
	}
}

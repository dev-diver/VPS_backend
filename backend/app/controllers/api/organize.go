package api

import (
	"strconv"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOrganizesHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, err := strconv.Atoi(c.Params("companyID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid companyID"})
		}

		orgMap, err := GetOrganizes(c, db, uint(companyID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		root, err := buildTree(orgMap)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(root)
	}
}

func GetOrganizes(c *fiber.Ctx, db *database.Database, companyID uint) (map[uint]*dto.OrganizeResponse, error) {
	var organizes []models.Organize
	if err := db.DB.Preload("Members").Where("company_id = ?", companyID).Find(&organizes).Error; err != nil {
		return nil, err
	}

	orgMap := make(map[uint]*dto.OrganizeResponse)
	for _, organize := range organizes {
		organizeDTO := dto.MapOrganizeToResponse(organize)
		orgMap[organizeDTO.ID] = &organizeDTO
	}

	return orgMap, nil
}

func buildTree(orgMap map[uint]*dto.OrganizeResponse) (dto.OrganizeResponse, error) {

	var root *dto.OrganizeResponse
	for _, tempOrg := range orgMap {
		if tempOrg.ParentID != nil {
			parent, ok := orgMap[*tempOrg.ParentID]
			if !ok {
				return dto.OrganizeResponse{}, fiber.NewError(fiber.StatusInternalServerError, "parent not found")
			}
			parent.Children = append(parent.Children, tempOrg)
		} else if root == nil {
			root = tempOrg
		} else {
			return dto.OrganizeResponse{}, fiber.NewError(fiber.StatusInternalServerError, "multiple root found")
		}
	}

	if root == nil {
		return dto.OrganizeResponse{}, fiber.NewError(fiber.StatusInternalServerError, "root not found")
	}

	return *root, nil
}

func AddOrganizeHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		companyID, err := strconv.ParseUint(c.Params("companyID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid companyID"})
		}

		organizeID, err := strconv.ParseUint(c.Params("organizeID"), 10, 32)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid companyID"})
		}

		organizeRequest := dto.OrganizeRequest{}
		if err := c.BodyParser(&organizeRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		parentID := uint(organizeID)

		organize := models.Organize{
			Name:      organizeRequest.Name,
			CompanyID: uint(companyID),
			ParentID:  &parentID,
		}

		if err := db.DB.Create(&organize).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateOrganizeHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		organizeID, err := strconv.ParseUint(c.Params("organizeID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid companyID"})
		}

		//요청 파싱
		organizeRequest := models.Organize{}
		if err := c.BodyParser(&organizeRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		//원래 값
		var organize models.Organize
		if err := db.DB.First(&organize, organizeID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "organize not found"})
		}

		//업데이트
		organize.Name = organizeRequest.Name
		if err := db.DB.Save(&organize).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func DeleteOrganizeHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		organizeID, err := strconv.ParseUint(c.Params("organizeID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organizeID"})
		}

		// 트랜잭션 시작
		tx := db.DB.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not start transaction"})
		}

		var organize models.Organize
		if err := tx.First(&organize, organizeID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "organize not found"})
		}

		// 모든 하위 조직 삭제
		if err := deleteSubOrganizes(tx, uint(organizeID)); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 해당 조직의 멤버 등록 해제
		if err := tx.Model(&organize).Association("Members").Clear(); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 조직 자체 삭제
		if err := tx.Delete(&organize, uint(organizeID)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 트랜잭션 커밋
		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not commit transaction"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// 모든 하위 조직을 재귀적으로 삭제하는 함수
func deleteSubOrganizes(tx *gorm.DB, parentID uint) error {
	var subOrganizes []models.Organize
	if err := tx.Where("parent_id = ?", parentID).Find(&subOrganizes).Error; err != nil {
		return err
	}

	for _, subOrganize := range subOrganizes {
		if err := deleteSubOrganizes(tx, subOrganize.ID); err != nil {
			return err
		}

		// 해당 조직 삭제
		if err := tx.Delete(&subOrganize).Error; err != nil {
			return err
		}
	}

	return nil
}

func UpdateOrganizeMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		organizeID, err := strconv.ParseUint(c.Params("organizeID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organizeID"})
		}

		var organize models.Organize
		if err := db.DB.First(&organize, organizeID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "organize not found"})
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

		if err := db.DB.Model(&organize).Association("Members").Replace(&members); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

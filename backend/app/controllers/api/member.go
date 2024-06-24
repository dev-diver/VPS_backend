package api

import (
	"strconv"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateCompanyMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// URL에서 companyID 가져오기
		companyID, err := strconv.ParseUint(c.Params("companyID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
		}

		// DTO 객체를 생성
		var memberDTOs []dto.CreateMemberRequest
		if err := c.BodyParser(&memberDTOs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 데이터 유효성 검사
		validate := validator.New()
		for _, memberDTO := range memberDTOs {
			if err := validate.Struct(memberDTO); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
		}

		// DTO 데이터를 실제 모델로 변환
		var members []models.Member
		for _, memberDTO := range memberDTOs {

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(memberDTO.Password), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not hash password"})
			}

			member := models.Member{
				CompanyID: uint(companyID),
				Name:      memberDTO.Name,
				Email:     memberDTO.Email,
				Password:  string(hashedPassword),
				HireDate:  memberDTO.HireDate,
				IsActive:  true, // 자동으로 True 설정
			}
			members = append(members, member)
		}

		// 데이터를 저장합니다.
		if err := db.DB.Create(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		//memberResponse 로 member를 변환
		var memberResponses []dto.MemberResponse
		for _, member := range members {
			memberResponses = append(memberResponses, dto.MemberResponse{
				ID:       member.ID,
				Name:     member.Name,
				Email:    member.Email,
				HireDate: member.HireDate,
				IsActive: member.IsActive,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(memberResponses)
	}
}

func SearchMembersHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyword := c.Query("keyword")
		companyID := c.Params("companyID")
		var members []dto.MemberResponse
		query := db.DB.Table("members").
			Select("id, name, email, hire_date, is_active").
			Where("company_id = ? AND (name LIKE ? OR email LIKE ?)", companyID, "%"+keyword+"%", "%"+keyword+"%")

		if err := query.Scan(&members).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(members)
	}
}

func GetMemberProfileHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("memberID")
		var member dto.MemberResponse
		if err := db.DB.Table("members").First(&member, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(member)
	}
}

func DeactivateMemberHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("memberID")
		var member dto.MemberResponse
		if err := db.DB.Table("members").First(&member, id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		member.IsActive = false
		if err := db.DB.Table("members").Save(&member).Error; err != nil {
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

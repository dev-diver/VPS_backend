package api

import (
	"time"

	"cywell.com/vacation-promotion/app/auth"
	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func MakeAdminHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		company := new(models.Company)
		if err := c.BodyParser(company); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		tx := db.DB.Begin()
		company, organize, err := CreateCompany(tx, *company)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		companyID := company.ID
		Password := "1234"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not hash password"})
		}

		member := models.Member{
			CompanyID:  uint(companyID),
			OrganizeID: &organize.ID,
			Name:       "관리자",
			Email:      "admin@" + company.Name + ".co.kr",
			Password:   string(hashedPassword),
			HireDate:   time.Now(),
			IsActive:   true, // 자동으로 True 설정
		}

		if err := tx.Create(&member).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func LoginHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginRequest dto.LoginRequest
		var loginResponse dto.LoginResponse

		// 요청 바디 파싱
		if err := c.BodyParser(&loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		//vailidation
		var validate = validator.New()
		if err := validate.Struct(loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		member, err := auth.GetCorrectMember(loginRequest, c, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		groupIDs := make([]uint, len(member.Groups))
		for i, group := range member.Groups {
			groupIDs[i] = group.ID
		}

		memberResponse := dto.MapMemberToDTO(&member)
		loginResponse.Member = memberResponse
		loginResponse.CompanyID = member.CompanyID
		loginResponse.GroupIDs = groupIDs

		auth.SetSessionAndToken(c, &loginResponse)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func LogoutHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth.DestorySessionAndToken(c)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

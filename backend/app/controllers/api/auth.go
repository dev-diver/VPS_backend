package api

import (
	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/app/utils"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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

		var member models.Member
		//db에서 사용자 정보 가져오기
		if err := db.DB.Preload("Groups").Where("email = ?", loginRequest.Email).First(&member).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// 비밀번호 검증
		if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(loginRequest.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		token, err := utils.GenerateJWT(member.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
		}

		groupIDs := make([]uint, len(member.Groups))
		for i, group := range member.Groups {
			groupIDs[i] = group.ID
		}

		memberResponse := dto.MapMemberToDTO(member)
		loginResponse.Member = memberResponse
		loginResponse.CompanyID = member.CompanyID
		loginResponse.GroupIDs = groupIDs
		loginResponse.Token = token

		return c.Status(fiber.StatusOK).JSON(loginResponse)
	}
}

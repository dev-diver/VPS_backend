package api

import (
	"cywell.com/vacation-promotion/app/auth"
	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/utils"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

		member, err := auth.GetCorrectMember(loginRequest, c, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		sessionStore, err := auth.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get session"})
		}
		defer sessionStore.Save()
		err = sessionStore.Regenerate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not regenerate session"})
		}

		groupIDs := make([]uint, len(member.Groups))
		for i, group := range member.Groups {
			groupIDs[i] = group.ID
		}

		memberResponse := dto.MapMemberToDTO(member)
		loginResponse.Member = memberResponse
		loginResponse.CompanyID = member.CompanyID
		loginResponse.GroupIDs = groupIDs

		token, err := utils.GenerateJWT(&loginResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "authToken_info",
			Value:    token,
			HTTPOnly: false,
			SameSite: "Lax",
		})

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func LogoutHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionStore, err := auth.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get session"})
		}

		err = sessionStore.Destroy()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not destroy session"})
		}

		c.ClearCookie("authToken_info")
		return c.SendStatus(fiber.StatusNoContent)
	}
}

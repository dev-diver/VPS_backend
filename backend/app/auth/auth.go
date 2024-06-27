package auth

import (
	"errors"
	"log"
	"time"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/app/utils"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SessionStore = session.New(session.Config{
	KeyLookup:      "cookie:authToken_session",
	Expiration:     24 * time.Hour,
	CookieHTTPOnly: true,
})

func AuthCheckMiddleware(c *fiber.Ctx) error {
	_, err := CheckToken(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.Next()
}

func GetCorrectMember(loginRequest dto.LoginRequest, c *fiber.Ctx, db *database.Database) (models.Member, error) {
	var member models.Member
	//db에서 사용자 정보 가져오기
	if err := db.DB.Preload("Groups").Where("email = ?", loginRequest.Email).First(&member).Error; err != nil {
		return models.Member{}, errors.New("invalid credentials")
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(loginRequest.Password)); err != nil {
		return models.Member{}, errors.New("invalid credentials")
	}
	return member, nil
}

func CheckToken(c *fiber.Ctx) (*dto.LoginResponse, error) {
	tokenString := c.Cookies("authToken_info")
	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecretKey(), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims.Auth, nil
}

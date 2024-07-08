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
	KeyLookup:         "cookie:authToken_session",
	CookieHTTPOnly:    true,
	CookieSessionOnly: true,
})

func AuthCheckMiddleware(c *fiber.Ctx) error {
	if err := CheckToken(c); err != nil {
		log.Println(err)
		DestorySessionAndToken(c)
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

func CheckToken(c *fiber.Ctx) error {

	tokenString := c.Cookies("authToken_info")
	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecretKey(), nil
	})
	if err != nil {
		return errors.New("cannot parse token")
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		return errors.New("invalid token")
	}

	tokenMemberID := claims.Auth.Member.ID
	// fmt.Printf("tokenMemberID: %v\n", tokenMemberID)

	// session token check
	session, err := SessionStore.Get(c)
	if err != nil {
		return errors.New("cannot get session")
	}

	member_id := session.Get("member_id")
	// fmt.Printf("member_id: %v\n", member_id)
	if member_id != tokenMemberID {
		return errors.New("invalid session token")
	}

	return nil
}

func SetSessionAndToken(c *fiber.Ctx, loginResponse *dto.LoginResponse) error {
	session, err := SessionStore.Get(c)
	if err != nil {
		return err
	}
	defer session.Save()
	err = session.Regenerate()
	if err != nil {
		return err
	}

	session.Set("member_id", loginResponse.Member.ID)

	token, err := utils.GenerateJWT(loginResponse)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:        "authToken_info",
		Value:       token,
		HTTPOnly:    false,
		SessionOnly: false,
		SameSite:    "Lax",
	})

	return nil
}

func DestorySessionAndToken(c *fiber.Ctx) {
	session, err := SessionStore.Get(c)
	if err != nil {
		return
	}

	err = session.Destroy()
	if err != nil {
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:        "authToken_info",
		HTTPOnly:    false,
		SessionOnly: false,
		Expires:     time.Now().Add(-(time.Hour * 2)),
		SameSite:    "Lax",
	})
}

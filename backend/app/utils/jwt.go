package utils

import (
	"fmt"
	"log"
	"os"

	"cywell.com/vacation-promotion/app/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte

type Config struct {
	JWTSecret string `json:"jwt_secret"`
}
type Claims struct {
	Auth *dto.LoginResponse `json:"auth"`
	jwt.RegisteredClaims
}

func GenerateJWT(authInfo *dto.LoginResponse) (string, error) {
	claims := &Claims{
		authInfo,
		jwt.RegisteredClaims{
			ID: fmt.Sprint(authInfo.Member.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func SetJWTSecretKey() error {

	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Error loading ./config/.env file: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	jwtKey = []byte(secret)
	return nil
}

func GetJWTSecretKey() []byte {
	return jwtKey
}

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

	err := godotenv.Load()
	if err != nil {
		log.Printf(".env 파일을 찾을 수 없습니다. 시스템 환경 변수를 사용합니다.")
	}

	secret := os.Getenv("JWT_SECRET")
	print("secret:", secret)
	if secret == "" {
		return fmt.Errorf("JWT_SECRET 환경 변수를 설정해주세요")
	}
	jwtKey = []byte(secret)
	return nil
}

func GetJWTSecretKey() []byte {
	return jwtKey
}

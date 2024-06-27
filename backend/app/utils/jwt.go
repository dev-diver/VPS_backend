package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"cywell.com/vacation-promotion/app/dto"
	"github.com/golang-jwt/jwt/v5"
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
			ID: string(authInfo.Member.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func SetJWTSecretKey() error {
	config := &Config{}
	file, err := os.ReadFile("config/secret.json")
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return fmt.Errorf("could not unmarshal config JSON: %w", err)
	}
	jwtKey = []byte(config.JWTSecret)
	return nil
}

func GetJWTSecretKey() []byte {
	return jwtKey
}

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

type Config struct {
	JWTSecret string `json:"jwt_secret"`
}
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
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

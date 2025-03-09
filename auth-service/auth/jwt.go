package auth

import (
	"time"

	"auth-service/config"
	"github.com/golang-jwt/jwt/v5"
)

// Генерация JWT
func GenerateJWT(uuid string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Токен на 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// Проверка JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
}

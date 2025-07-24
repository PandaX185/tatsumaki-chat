package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(claims map[string]string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   claims["userId"],
		"username": claims["username"],
		"fullname": claims["fullname"],
	}).SignedString([]byte(os.Getenv("AUTH_KEY")))
}

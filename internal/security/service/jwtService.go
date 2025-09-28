package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = email
	claims["exp"] = time.Now().Add(time.Hour * 192).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "Failed to Create jwt token", err
	}
	return tokenString, nil
}

package xjwt

import "github.com/golang-jwt/jwt/v5"

func GenerateJWT() (string, error) {

	jwt.New(jwt.SigningMethodHS512, func(t *jwt.Token) {})
	return "", nil
}

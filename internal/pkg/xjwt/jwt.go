package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTHelper interface {
	GenerateJWT(issuer string, userID string, expiryTime time.Duration, key string) (string, error)
}

type XJWT struct{}

func (x *XJWT) GenerateJWT(issuer string, userID string, expiryTime time.Duration, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Issuer:    issuer,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

package xjwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTHelper interface {
	GenerateJWT(issuer string, userID string, expiryTime time.Duration, key string) (string, error)
	ValidateJWT(token string, secret string) (*jwt.Token, error)
	GetTokenClaims(token string, secret string) (jwt.Claims, error)
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

func (x *XJWT) ValidateJWT(token string, secret string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return parsedToken, fmt.Errorf("invalid token")
	}

	if _, ok := parsedToken.Claims.(jwt.MapClaims); !ok {
		return parsedToken, fmt.Errorf("invalid jwt claims")
	}

	return parsedToken, nil
}

func (x *XJWT) GetTokenClaims(token string, secret string) (jwt.Claims, error) {
	parsedToken, err := x.ValidateJWT(token, secret)
	if err != nil {
		return jwt.MapClaims{}, err
	}

	return parsedToken.Claims, nil
}

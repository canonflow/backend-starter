package pkg

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var TOKEN_COOKIE = "Authorization"

type AccessTokenContextWrapper interface {
	Context() context.Context
	GetToken() (string, error)
	SetToken(value string, durationInMinutes uint)
	DeleteToken()
}

type (
	userKey struct{}
)

func CreateAccessToken(id int, email, key string, durationInMinutes uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        id,
		"email":      email,
		"iat":        jwt.NewNumericDate(now),
		"expired_at": jwt.NewNumericDate(now.Add(time.Minute * time.Duration(durationInMinutes))),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Cannot claim JWT Token")
	}

	return claims, nil
}

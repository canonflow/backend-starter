package pkg

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	TOKEN_COOKIE    = "Authorization"
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrMissingClaim = errors.New("required claim missing or malformed")
	ErrTokenExpired = errors.New("token has expired")

	JWTUserIDKey = "JWT_USER_ID_KEY"
	JWTEmailKey = "JWT_EMAIL_KEU"
)

type AccessTokenContextWrapper interface {
	Context() context.Context
	GetToken() (string, error)
	SetToken(value string, durationInMinutes uint)
	DeleteToken()
	SetLocal(key string, value any)
	Local(key string) any
}

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
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot claim JWT token")
	}

	// Check expired_at BEFORE trusting/extracting anything else from claims
	if err := checkExpiration(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func GetUserIDFromClaims(claims jwt.MapClaims) (int, error) {
	raw, ok := claims["sub"]
	if !ok {
		return 0, ErrMissingClaim
	}

	f, ok := raw.(float64)
	if !ok {
		return 0, ErrMissingClaim
	}

	return int(f), nil
}

func GetEmailFromClaims(claims jwt.MapClaims) (string, error) {
	raw, ok := claims["email"]
	if !ok {
		return "", ErrMissingClaim
	}

	email, ok := raw.(string)
	if !ok {
		return "", ErrMissingClaim
	}

	return email, nil
}

// checkExpiration validates the custom "expired_at" claim manually,
// since it's not the standard "exp" claim jwt/v5 auto-validates.
func checkExpiration(claims jwt.MapClaims) error {
	raw, ok := claims["expired_at"]
	if !ok {
		return ErrMissingClaim
	}

	expiredAtFloat, ok := raw.(float64) // JSON numbers unmarshal as float64
	if !ok {
		return ErrMissingClaim
	}

	expiredAt := time.Unix(int64(expiredAtFloat), 0)

	if time.Now().After(expiredAt) {
		return ErrTokenExpired
	}

	return nil
}

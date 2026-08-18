package pkg

import (
	"github.com/canonflow/backend-starter/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), config.Get[int](config.BcryptCost))
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CheckHash(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}

package repository_test

import (
	"testing"

	"github.com/canonflow/backend-starter/internal/app/security/repository"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt_Valid(t *testing.T) {
	security := repository.NewAESCBCSecurity()
	key := "giY9sbzrjCj6PMEFVE5VNBhKesbzdu5N1RRAZbHN6io=" // pragma: allow secrets
	message := "test@example.com"

	encrypted, err := security.Encrypt(message, key)

	assert.Equal(t, err, nil)
	t.Log("Encrypted:", encrypted)

	decrypted, err := security.Decrypt(encrypted, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, message, decrypted)
}

package repository

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type aesCbcSecurity struct{}

func NewAESCBCSecurity() ISecurity {
	return &aesCbcSecurity{}
}

// KeySize for AES-256 is 32 bytes.
const KeySize = 32

var (
	ErrInvalidKeySize = errors.New("key must be exactly 32 bytes for AES-256")
	ErrInvalidCipher  = errors.New("ciphertext is too short or not a multiple of the block size")
	ErrInvalidPadding = errors.New("invalid padding")
)

// Encrypt encrypts plaintext with AES-256-CBC.
// A random IV is generated and prepended to the ciphertext.
// The result is returned as a base64 string: base64(IV || ciphertext).
func (r *aesCbcSecurity) Encrypt(plainText, key string) (string, error) {
	k, err := decodeKey(key)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(k)
	plainTextBytes := []byte(plainText)

	if len(keyBytes) != KeySize {
		return "", ErrInvalidKeySize
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Copy so the caller's slice is never modified by append in pkcs7Pad
	padded := pkcs7Pad(append([]byte(nil), plainTextBytes...), block.BlockSize())

	out := make([]byte, aes.BlockSize+len(padded))
	iv := out[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCBCEncrypter(block, iv).
		CryptBlocks(out[aes.BlockSize:], padded)

	return base64.StdEncoding.EncodeToString(out), nil
}

// Decrypt reverses Encrypt: it takes base64(IV || ciphertext) and returns the plaintext.
func (r *aesCbcSecurity) Decrypt(cipherText, key string) (string, error) {
	k, err := decodeKey(key)
	if err != nil {
		return "", err
	}
	keyBytes := []byte(k)

	if len(keyBytes) != KeySize {
		return "", ErrInvalidKeySize
	}

	raw, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(raw) < 2*aes.BlockSize || len(raw)%aes.BlockSize != 0 {
		return "", ErrInvalidCipher
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	iv := raw[:aes.BlockSize]
	ciphertext := raw[aes.BlockSize:]

	plain := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, ciphertext)

	result, err := pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// pkcs7Pad pads data to a multiple of blockSize.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

// pkcs7Unpad removes PKCS#7 padding.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, ErrInvalidPadding
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(data) {
		return nil, ErrInvalidPadding
	}
	for _, b := range data[len(data)-padLen:] {
		if int(b) != padLen {
			return nil, ErrInvalidPadding
		}
	}
	return data[:len(data)-padLen], nil
}

func decodeKey(key string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("key is not valid base64: %w", err)
	}
	if len(raw) != KeySize {
		return nil, ErrInvalidKeySize
	}
	return raw, nil
}

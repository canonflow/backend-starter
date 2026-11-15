package repository

type ISecurity interface {
	Encrypt(text, key string) (string, error)
	Decrypt(cipherText, key string) (string, error)
}

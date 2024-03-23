package utils

import (
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordSalt() ([]byte, error) {
	salt := make([]byte, 16) // TODO more bytes?
	_, err := rand.Read(salt)
	return salt, err
}

func HashPassword(password string, salt []byte) (*uint, error) {
	salted := password + string(salt)
	bytes, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(bytes)
	var result uint = 0
	for i, b := range hashed {
		shift := uint(i) % 8
		result = result + uint(b)<<shift
	}
	return &result, nil
}

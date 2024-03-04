package utils

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordSalt() ([]byte, error) {
	salt := make([]byte, 16) // TODO more bytes?
	_, err := rand.Read(salt)
	return salt, err
}

func HashPassword(password string, salt []byte) ([]byte, error) {
	salted := password + string(salt)
	return bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
}

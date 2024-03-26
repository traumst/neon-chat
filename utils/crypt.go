package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strconv"

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
	var hash uint
	if strconv.IntSize == 64 {
		fold := Fold(bytes)
		hash = uint(binary.BigEndian.Uint64(fold[:]))
	} else {
		return nil, fmt.Errorf("strconv.IntSize expected to be 64, but is [%d]", strconv.IntSize)
	}
	return &hash, nil
}

func Fold(bytes []byte) [8]byte {
	var fold [8]byte
	for i := 0; i < len(bytes); i += 8 {
		for j := 0; j < 8; j++ {
			if i+j < len(bytes) {
				fold[j] ^= bytes[i+j]
			}
		}
	}
	return fold
}

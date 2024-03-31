package utils

import (
	"log"
)

func HashPassword(password string, salt string) (string, error) {
	s := string(salt[:])
	log.Printf("-----> HashPassword TRACE IN pass[%s] salt[%s]\n", password, s)
	salted := password + s
	bytes := []byte(salted)
	fold := FoldInto8Bytes(bytes)
	return string(fold[:]), nil
}

func FoldInto8Bytes(bytes []byte) [8]byte {
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

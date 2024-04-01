package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

// TODO think: type forces salt change when switching user.type
func GenerateSalt(userName string, userType string) string {
	seed := fmt.Sprintf("%s-%s", userType, userName)
	saltPlain := fmt.Sprintf("%s;%s", RandStringBytes(7), seed)
	salt := ToHexSha256(saltPlain)
	if salt == "" || strings.Contains(salt, "\n") || strings.Contains(salt, " ") {
		log.Fatalf("failed to generate salt for user[%s] generated[%s] from [%s]", userName, salt, seed)
	}
	return salt
}

func HashPassword(password string, salt string) (string, error) {
	salted := password + salt
	hash := ToHexSha256(salted)
	return hash, nil
}

func ToHexSha256(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

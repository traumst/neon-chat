package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func GenerateSalt(userName string, userType string) string {
	rr := int(rand.Int31n(8)) + 3
	saltPlain := fmt.Sprintf("%s%s%s%s%s",
		RandStringBytes(rr),
		userName,
		RandStringBytes(rr),
		userType,
		RandStringBytes(rr))
	salt := ToHexSha256(saltPlain)
	if salt == "" || strings.Contains(salt, "\n") || strings.Contains(salt, " ") {
		log.Fatalf("failed to generate salt for user[%s] generated salt[%s]", userName, salt)
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

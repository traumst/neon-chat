package utils

import (
	"math/rand"
	"neon-chat/src/consts"
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	len := len(consts.LetterBytes)
	for i := range b {
		idx := rand.Intn(len)
		b[i] = consts.LetterBytes[idx]
	}
	return string(b)
}

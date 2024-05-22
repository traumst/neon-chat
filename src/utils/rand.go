package utils

import "math/rand"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	len := len(LetterBytes)
	for i := range b {
		idx := rand.Intn(len)
		b[i] = LetterBytes[idx]
	}
	return string(b)
}

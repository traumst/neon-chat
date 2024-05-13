package utils

import (
	"fmt"
	"log"
)

const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

func Size(bytes int64) string {
	log.Printf("formatting size: %#v", bytes)
	switch {
	case bytes < KB:
		return fmt.Sprintf("%dBytes", bytes)
	case bytes < MB:
		return fmt.Sprintf("%dKB", bytes/KB)
	case bytes < GB:
		return fmt.Sprintf("%dMB", bytes/MB)
	default:
		panic(fmt.Sprintf("size is over 1GB: %#v", bytes))
	}
}

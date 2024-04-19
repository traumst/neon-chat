package utils

import (
	"fmt"
	"log"
)

func Size(bytes int64) string {
	log.Printf("formatting size: %#v", bytes)
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dBytes", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%dKB", bytes/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%dMB", bytes/1024*1024)
	default:
		panic(fmt.Sprintf("size is over 1GB: %#v", bytes))
	}
}

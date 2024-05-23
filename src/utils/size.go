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

func SizeEncode(bytes int64) string {
	log.Printf("formatting size: %#v", bytes)
	switch {
	case bytes < KB:
		return fmt.Sprintf("%dBytes", bytes)
	case bytes < MB:
		return fmt.Sprintf("%dKB", bytes/KB)
	case bytes < GB:
		return fmt.Sprintf("%dMB", bytes/MB)
	default:
		panic(fmt.Sprintf("size is over 1GB: [%d]", bytes))
	}
}

func SizeDecode(str string) int64 {
	var bytes int64
	var unit string
	_, err := fmt.Sscanf(str, "%d%s", &bytes, &unit)
	if err != nil {
		panic(fmt.Errorf("failed decoding [%s], %s", str, err.Error()))
	}
	switch unit {
	case "Bytes":
		return bytes
	case "KB":
		return bytes * KB
	case "MB":
		return bytes * MB
	default:
		panic(fmt.Sprintf("size is over 1GB: [%s]", str))
	}
}

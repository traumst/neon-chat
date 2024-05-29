package utils

import (
	"fmt"
)

const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

func SizeEncode(bytes int64) string {
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

func SizeDecode(size string) int64 {
	var bytes int64
	var unit string
	_, err := fmt.Sscanf(size, "%d%s", &bytes, &unit)
	if err != nil {
		panic(fmt.Errorf("failed decoding size[%s], %s", size, err.Error()))
	}
	switch unit {
	case "Bytes":
		return bytes
	case "KB":
		return bytes * KB
	case "MB":
		return bytes * MB
	default:
		panic(fmt.Sprintf("size is over 1GB: [%s]", size))
	}
}

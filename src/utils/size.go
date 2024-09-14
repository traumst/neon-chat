package utils

import (
	"fmt"
	"neon-chat/src/consts"
)

func SizeEncode(bytes int64) string {
	switch {
	case bytes < consts.KB:
		return fmt.Sprintf("%dBytes", bytes)
	case bytes < consts.MB:
		return fmt.Sprintf("%dKB", bytes/consts.KB)
	case bytes < consts.GB:
		return fmt.Sprintf("%dMB", bytes/consts.MB)
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
		return bytes * consts.KB
	case "MB":
		return bytes * consts.MB
	default:
		panic(fmt.Sprintf("size is over 1GB: [%s]", size))
	}
}

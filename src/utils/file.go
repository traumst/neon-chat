package utils

import (
	"os"
)

func GetFilenamesIn(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	list := make([]string, 0)
	for _, entry := range entries {
		list = append(list, entry.Name())
	}
	return list, nil
}

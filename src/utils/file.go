package utils

import (
	"log"
	"os"
)

func GetFilenamesIn(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("GetFilenamesIn error getting current directory:", err)
	}
	list := make([]string, 0)
	for _, entry := range entries {
		list = append(list, entry.Name())
	}
	return list, nil
}

func LS() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("LS error getting current directory:", err)
		return
	} else {
		log.Println("LS current directory:", dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println("LS error reading directory contents:", err)
		return
	} else {
		log.Println("LS content of", dir, "is", files)
	}
}

package utils

import (
	"fmt"
	"io"
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

func ReadFileContent(filePath string) (string, error) {
	fileContent, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open .env file from [%s]: %s", filePath, err)
	}
	defer fileContent.Close()

	contentBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to read content from file[%s]: %s", filePath, err)
	}
	content := string(contentBytes)
	return content, nil
}

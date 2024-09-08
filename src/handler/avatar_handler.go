package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	a "neon-chat/src/model/app"
	"neon-chat/src/utils"
	"net/http"
)

var allowedImageFormats = []string{
	"image/svg+xml",
	"image/jpeg",
	"image/gif",
	"image/png",
}

func isAllowedImageFormat(mime string) bool {
	for _, allowed := range allowedImageFormats {
		if allowed == mime {
			return true
		}
	}
	return false
}

func UpdateAvatar(
	db *d.DBConn,
	userId uint,
	file *multipart.File,
	info *multipart.FileHeader,
) (*a.Avatar, error) {
	if info.Size > utils.MaxUploadBytesSize {
		return nil, fmt.Errorf("file too large %d, limit is %d", info.Size, utils.MaxUploadBytesSize)
	} else if len(info.Filename) == 0 {
		return nil, fmt.Errorf("file lacks name")
	} else if len(info.Filename) > utils.MaxFileName {
		return nil, fmt.Errorf("file name is too long")
	}
	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file[%s]", info.Filename)
	}
	mime := http.DetectContentType(fileBytes)
	if !isAllowedImageFormat(mime) {
		return nil, fmt.Errorf("file type is not supported[%s]", mime)
	}

	return shared.UpdateAvatar(db, userId, info.Filename, fileBytes, mime)
}

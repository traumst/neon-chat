package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	d "prplchat/src/db"
	"prplchat/src/utils"
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
) (*d.Avatar, error) {
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
	oldAvatars, err := db.GetUserAvatars(userId)
	if err != nil {
		return nil, fmt.Errorf("fail to get avatar for user[%d]", userId)
	}
	saved, err := db.AddAvatar(userId, info.Filename, fileBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar[%s], %s", info.Filename, err.Error())
	}
	if len(oldAvatars) > 0 {
		for _, old := range oldAvatars {
			if old == nil {
				continue
			}
			err := db.DropAvatar(old.Id)
			if err != nil {
				log.Printf("updateAvatar ERROR failed to drop old avatar[%v]", old)
			}
		}
	}
	return saved, nil
}

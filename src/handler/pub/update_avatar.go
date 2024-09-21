package pub

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/jmoiron/sqlx"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/utils"
)

func UpdateAvatar(
	dbConn *db.DBConn,
	userId uint,
	file *multipart.File,
	info *multipart.FileHeader,
) (*app.Avatar, error) {
	if info.Size > consts.MaxUploadBytesSize {
		return nil, fmt.Errorf("file too large %d, limit is %d", info.Size, consts.MaxUploadBytesSize)
	} else if len(info.Filename) == 0 {
		return nil, fmt.Errorf("file lacks name")
	} else if len(info.Filename) > consts.MaxFileName {
		return nil, fmt.Errorf("file name is too long")
	}
	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file[%s]", info.Filename)
	}
	mime := http.DetectContentType(fileBytes)
	if !utils.IsAllowedImageFormat(mime) {
		return nil, fmt.Errorf("file type is not supported[%s]", mime)
	}

	return setAvatar(dbConn.Tx, userId, info.Filename, fileBytes, mime)
}

func setAvatar(
	dbConn sqlx.Ext,
	userId uint,
	filename string,
	fileBytes []byte,
	mime string,
) (*app.Avatar, error) {
	oldAvatars, err := db.GetUserAvatars(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("fail to get avatar for user[%d]", userId)
	}
	saved, err := db.AddAvatar(dbConn, userId, filename, fileBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar[%s], %s", filename, err.Error())
	}
	if len(oldAvatars) > 0 {
		for _, old := range oldAvatars {
			if old == nil {
				continue
			}
			err := db.DropAvatar(dbConn, old.Id)
			if err != nil {
				log.Printf("ERROR setAvatar failed to drop old avatar[%v]", old)
			}
		}
	}
	return convert.AvatarDBToApp(saved), nil
}

package convert

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
	"prplchat/src/utils"
)

func AvatarAppToDB(avatar *app.Avatar) *db.Avatar {
	return &db.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   int(utils.SizeDecode(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}

func AvatarDBToApp(avatar *db.Avatar) *app.Avatar {
	return &app.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   utils.SizeEncode(int64(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}
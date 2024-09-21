package convert

import (
	"neon-chat/src/app"
	"neon-chat/src/db"
)

func ChatAppToDB(chat *app.Chat) *db.Chat {
	if chat == nil {
		return nil
	}
	return &db.Chat{
		Id:      chat.Id,
		Title:   chat.Name,
		OwnerId: chat.OwnerId,
	}
}

func ChatDBToApp(chat *db.Chat, owner *db.User) *app.Chat {
	if chat == nil {
		return nil
	}
	return &app.Chat{
		Id:        chat.Id,
		Name:      chat.Title,
		OwnerId:   owner.Id,
		OwnerName: owner.Name,
	}
}

package convert

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
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

func ChatDBToApp(chat *db.Chat) *app.Chat {
	if chat == nil {
		return nil
	}
	return &app.Chat{
		Id:      chat.Id,
		Name:    chat.Title,
		OwnerId: chat.OwnerId,
	}
}

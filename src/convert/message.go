package convert

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
)

func MessageAppToDB(message *app.Message) db.Message {
	return db.Message{
		Id:       message.Id,
		ChatId:   message.ChatId,
		AuthorId: message.Author.Id,
		Text:     message.Text,
	}
}

func MessageDBToApp(message *db.Message, author *app.User) app.Message {
	return app.Message{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: author,
		Text:   message.Text,
	}
}

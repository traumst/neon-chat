package convert

import (
	"neon-chat/src/app"
	"neon-chat/src/db"
)

func MessageAppToDB(message *app.Message) db.Message {
	return db.Message{
		Id:       message.Id,
		ChatId:   message.ChatId,
		AuthorId: message.Author.Id,
		Text:     message.Text,
	}
}

func MessageDBToApp(message *db.Message, author *app.User, quote *app.Message) app.Message {
	return app.Message{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: author,
		Text:   message.Text,
		Quote:  quote,
	}
}

func MessageDBToQuoteApp(message *db.Message, author *app.User) app.Quote {
	return app.Quote{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: author,
		Text:   message.Text,
	}
}

func MessageAppToQuoteApp(message *app.Message) app.Quote {
	return app.Quote{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: message.Author,
		Text:   message.Text,
	}
}

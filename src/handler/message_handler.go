package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/sse"
	"neon-chat/src/handler/state"
	i "neon-chat/src/interface"
	a "neon-chat/src/model/app"
	"neon-chat/src/model/event"
	t "neon-chat/src/model/template"
)

// TODO consider adding quote thread with depth limit
func HandleGetQuote(
	state *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	msgId uint,
) (string, error) {
	log.Printf("HandleGetQuote TRACE quoting message[%d] of chat[%d]\n", msgId, chatId)
	canChat, err := d.UsersCanChat(db.Conn, chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetQuote ERROR checking whether user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return "", fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetQuote ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return "", fmt.Errorf("user is not in chat")
	}
	//
	dbOwner, err := d.GetOwner(db.Conn, chatId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting owner for chat[%d], %s\n", chatId, err.Error())
		return "", fmt.Errorf("failed to get chat owner from db, %s", err.Error())
	}
	appOwner := convert.UserDBToApp(dbOwner, nil)
	//
	dbMsg, err := d.GetMessage(db.Conn, msgId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting quote[%d] from db, %s\n", msgId, err)
		return "", fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	appQuote := convert.MessageDBToQuoteApp(dbMsg, user)
	//
	dbAvatar, err := d.GetAvatar(db.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return "nil", fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}
	//
	dbAuthor, err := d.GetUser(db.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] from db, %s\n", dbMsg.AuthorId, err)
		return "", fmt.Errorf("failed to get author from db: %s", err.Error())
	}
	appQuote.Author = convert.UserDBToApp(dbAuthor, dbAvatar)
	//
	msgTmpl, err := appQuote.Template(user, appOwner)
	if err != nil {
		log.Printf("HandleGetQuote ERROR generating message template, %s\n", err)
		return "nil", fmt.Errorf("failed to generate message template: %s", err.Error())
	}
	quoteTmpl := &t.QuoteTemplate{
		IntermediateId: msgTmpl.IntermediateId,
		ChatId:         msgTmpl.ChatId,
		MsgId:          msgTmpl.MsgId,
		AuthorId:       msgTmpl.AuthorId,
		AuthorName:     msgTmpl.AuthorName,
		AuthorAvatar:   msgTmpl.AuthorAvatar,
		Text:           msgTmpl.Text,
		TextIntro:      msgTmpl.TextIntro,
	}
	return quoteTmpl.HTML()
}

func HandleGetMessage(
	state *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	msgId uint,
) (i.Renderable, error) {
	log.Printf("HandleGetMessage TRACE getting message[%d] from chat[%d]\n", msgId, chatId)
	canChat, err := d.UsersCanChat(db.Conn, chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetMessage ERROR checking whether user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return nil, fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetMessage ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return nil, fmt.Errorf("user is not in chat")
	}
	dbOwner, err := d.GetOwner(db.Conn, chatId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting owner for chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to get message from db, %s", err.Error())
	}
	dbMsg, err := d.GetMessage(db.Conn, msgId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting message[%d] from db, %s\n", msgId, err)
		return nil, fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	dbAvatar, err := d.GetAvatar(db.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}
	var appQuote *a.Message
	if dbQuote, _ := d.GetQuote(db.Conn, msgId); dbQuote != nil {
		dbMsg, err := d.GetMessage(db.Conn, dbQuote.QuoteId)
		if err != nil {
			log.Printf("HandleGetMessage warn getting quote message[%d] from db, %s\n", dbQuote.QuoteId, err)
			//return nil, fmt.Errorf("failed to get quote message[%d] from db: %s", dbQuote.QuoteId, err.Error())
		}
		dbAvatar, err := d.GetAvatar(db.Conn, dbMsg.AuthorId)
		if err != nil {
			log.Printf("HandleGetMessage ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
			return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
		}
		user.Avatar = convert.AvatarDBToApp(dbAvatar)
		quoteMsg := convert.MessageDBToApp(dbMsg, user, nil)
		appQuote = &quoteMsg
	}
	appOwner := convert.UserDBToApp(dbOwner, nil)
	appMsg := convert.MessageDBToApp(dbMsg, user, appQuote)
	appMsg.Author.Avatar = convert.AvatarDBToApp(dbAvatar)
	tmplMsg, err := appMsg.Template(user, appOwner, appQuote)
	if err != nil {
		log.Printf("HandleGetMessage ERROR generating message template, %s\n", err)
		return nil, fmt.Errorf("failed to generate message template: %s", err.Error())
	}
	return &tmplMsg, nil
}

// TODO has partial success state - if fail of saving quote
func HandleMessageAdd(
	state *state.State,
	db *d.DBConn,
	chatId uint,
	author *a.User,
	msg string,
	quoteId uint,
) (*a.Message, error) {
	log.Printf("HandleMessageAdd TRACE opening current chat for user[%d]\n", author.Id)
	canChat, err := d.UsersCanChat(db.Tx, chatId, author.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user[%d] can chat[%d]: %s", author.Id, chatId, err.Error())
	}
	if !canChat {
		return nil, fmt.Errorf("user is not in chat")
	}
	dbChat, err := d.GetChat(db.Tx, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	dbMsg, err := d.AddMessage(db.Tx, &d.Message{
		Id:       0,
		ChatId:   chatId,
		AuthorId: author.Id,
		Text:     msg,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	// quoteId 0 means message has no quote attached
	var appQuote *a.Message
	if quoteId != 0 {
		quote, err := d.AddQuote(db.Tx, &d.Quote{
			MsgId:   dbMsg.Id,
			QuoteId: quoteId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add quote[%d] to message[%d]: %s", quoteId, dbMsg.Id, err.Error())
		}
		dbQuote, err := d.GetMessage(db.Tx, quote.QuoteId)
		if err != nil {
			return nil, fmt.Errorf("failed to get quote[%d] from db: %s", quoteId, err.Error())
		}
		quoteAuthor, err := shared.GetUser(db.Tx, dbQuote.AuthorId)
		if err != nil {
			return nil, fmt.Errorf("failed to get quote[%d] author[%d] avatar from db: %s", quoteId, dbQuote.AuthorId, err.Error())
		}
		tmp := convert.MessageDBToApp(dbQuote, quoteAuthor, nil)
		appQuote = &tmp
	}
	dbOwner, err := d.GetUser(db.Tx, dbChat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", chatId, dbChat.OwnerId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat, dbOwner)
	appMsg := convert.MessageDBToApp(dbMsg, author, appQuote)
	err = sse.DistributeMsg(state, db.Tx, appChat, &appMsg, event.MessageAdd)
	if err != nil {
		log.Printf("HandleMessageAdd ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, nil
}

func HandleMessageDelete(
	state *state.State,
	db *d.DBConn,
	chatId uint,
	user *a.User,
	msgId uint,
) (*a.Message, error) {
	log.Printf("HandleMessageDelete TRACE removing msg[%d] from chat[%d] for user[%d]\n", msgId, chatId, user.Id)
	dbMsg, err := d.GetMessage(db.Tx, msgId)
	if err != nil {
		return nil, fmt.Errorf("failed to get message[%d] from db: %s", msgId, err.Error())
	}
	dbChat, err := d.GetChat(db.Tx, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	if dbMsg.AuthorId != user.Id && dbChat.OwnerId != user.Id {
		return nil, fmt.Errorf("user[%d] is not allowed to delete message[%d] in chat[%d]", user.Id, msgId, chatId)
	}
	err = d.DeleteMessage(db.Tx, msgId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove message[%d] from chat[%d] in db, %s", msgId, chatId, err.Error())
	}
	dbSpecialUsers, err := d.GetUsers(db.Tx, []uint{dbMsg.AuthorId, dbChat.OwnerId})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", dbChat.Id, dbChat.OwnerId, err.Error())
	}
	var appChatOwner *a.User
	var appMsgAuthor *a.User
	for _, dbSpecialUser := range dbSpecialUsers {
		if dbSpecialUser.Id == dbChat.OwnerId {
			appChatOwner = convert.UserDBToApp(&dbSpecialUser, nil)
		}
		if dbSpecialUser.Id == dbMsg.AuthorId {
			appMsgAuthor = convert.UserDBToApp(&dbSpecialUser, nil)
		}
	}
	if appChatOwner == nil {
		return nil, fmt.Errorf("chat[%d] owner[%d] not found", dbChat.Id, dbChat.OwnerId)
	}
	if appMsgAuthor == nil {
		return nil, fmt.Errorf("message[%d] author[%d] not found", dbMsg.Id, dbMsg.AuthorId)
	}
	appChat := convert.ChatDBToApp(dbChat, &d.User{
		Id:     appChatOwner.Id,
		Name:   appChatOwner.Name,
		Email:  appChatOwner.Email,
		Type:   string(appChatOwner.Type),
		Status: string(appChatOwner.Status),
		Salt:   appChatOwner.Salt,
	})
	if appChat == nil {
		return nil, fmt.Errorf("cannot convert chat[%d] for app, owner[%v]", dbChat.Id, appChatOwner)
	}
	appMsg := convert.MessageDBToApp(dbMsg, appMsgAuthor, nil) // TODO bad user
	err = sse.DistributeMsg(state, db.Tx, appChat, &appMsg, event.MessageDrop)
	if err != nil {
		log.Printf("HandleMessageDelete ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, err
}

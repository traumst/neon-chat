package sse

import (
	"fmt"
	"log"
	"sync"

	"prplchat/src/convert"
	"prplchat/src/db"
	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/event"
)

// targetUser=nil means all users in chat
func DistributeUserChange(
	state *state.State,
	db *db.DBConn,
	targetUser *app.User, // who to inform, nil for all users
	subjectUser *app.User, // which user changed
	updateType event.EventType,
) error {
	if subjectUser == nil {
		return fmt.Errorf("subject user is nil")
	}
	if targetUser != nil {
		return distributeUpdateOfUser(state, targetUser, subjectUser, updateType)
	}
	dbChats, err := db.GetUserChats(subjectUser.Id)
	if err != nil {
		return fmt.Errorf("failed to get chats for user[%d], %s", subjectUser.Id, err)
	}
	var chats []*app.Chat
	for _, dbChat := range dbChats {
		chats = append(chats, convert.ChatDBToApp(&dbChat))
	}
	if len(chats) <= 0 {
		log.Printf("userChanged WARN user[%d] has no chats\n", subjectUser.Id)
		return fmt.Errorf("user[%d] has no chats", subjectUser.Id)
	}
	var wg sync.WaitGroup
	for _, chat := range chats {
		if chat == nil {
			continue
		}
		wg.Add(1)
		go func(chat *app.Chat) {
			defer wg.Done()
			err := distributeInCommonChat(db, chat, state, subjectUser, updateType)
			if err != nil {
				log.Printf("userChanged ERROR failed to distribute to chat[%d], %s\n", chat.Id, err)
			}
		}(chat)
	}
	wg.Wait()
	return err
}

func distributeInCommonChat(
	db *db.DBConn,
	chat *app.Chat,
	state *state.State,
	subjectUser *app.User,
	updateType event.EventType,
) error {
	targetUsers, err := chat.GetUsers(subjectUser.Id)
	if err != nil {
		return fmt.Errorf("failed to get users in chat[%d] for subject[%d], %s", chat.Id, subjectUser.Id, err)
	}
	// inform if chat is open
	for _, targetUser := range targetUsers {
		if targetUser == nil {
			continue
		}
		openChatId := state.GetOpenChat(targetUser.Id)
		if openChatId == 0 || openChatId != chat.Id {
			continue
		}
		uErr := distributeUpdateOfUser(state, targetUser, subjectUser, updateType)
		if uErr != nil {
			err = fmt.Errorf("%s, %s", uErr.Error(), err.Error())
		}
	}
	return err
}

func distributeUpdateOfUser(
	state *state.State,
	targetUser *app.User,
	subjectUser *app.User,
	updateType event.EventType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DUOU FATAL: %v", r)
		}
	}()
	conns := state.GetConn(targetUser.Id)
	for _, conn := range conns {
		if conn.User.Id != targetUser.Id {
			return fmt.Errorf("user[%d] does not own conn[%v], user[%d] does", targetUser.Id, conn.Origin, conn.User.Id)
		}
		var connerr error
		switch updateType {
		case event.UserChange:
			connerr = userNameChanged(conn, subjectUser)
		default:
			connerr = fmt.Errorf("unknown event type[%v]", updateType)
		}
		if err == nil {
			err = connerr
		} else {
			err = fmt.Errorf("%s, %s", err.Error(), connerr.Error())
		}
	}
	return err
}

func userNameChanged(conn *state.Conn, subject *app.User) error {
	if subject == nil {
		return fmt.Errorf("subjectUser is nil for userChanged")
	}
	log.Printf("userChanged TRACE informing target[%d] about subject[%d] change\n", conn.User.Id, subject.Id)
	tmpl := subject.Template(0, 0, conn.User.Id)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template user")
	}
	conn.In <- event.LiveEvent{
		Event:    event.UserChange,
		ChatId:   0,
		UserId:   subject.Id,
		MsgId:    0,
		AuthorId: subject.Id,
		Data:     data,
	}
	return nil
}

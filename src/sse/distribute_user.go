package sse

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"neon-chat/src/handler/pub"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/state"
)

// targetUser=nil means all users in chat
func DistributeUserChange(
	state *state.State,
	dbConn sqlx.Ext,
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
	appChats, err := pub.GetChats(dbConn, subjectUser.Id)
	if err != nil {
		return fmt.Errorf("failed to get chats for user[%d], %s", subjectUser.Id, err)
	}
	if len(appChats) <= 0 {
		log.Printf("userChanged WARN user[%d] has no chats\n", subjectUser.Id)
		return fmt.Errorf("user[%d] has no chats", subjectUser.Id)
	}
	for _, chat := range appChats {
		if chat == nil {
			continue
		}
		err := distributeInCommonChat(dbConn, chat, state, subjectUser, updateType)
		if err != nil {
			log.Printf("userChanged ERROR failed to distribute to chat[%d], %s\n", chat.Id, err)
		}
	}
	return err
}

func distributeInCommonChat(
	dbConn sqlx.Ext,
	chat *app.Chat,
	state *state.State,
	subjectUser *app.User,
	updateType event.EventType,
) error {
	// TODO is bad bad
	targetUsers, err := pub.GetChatUsers(dbConn, chat.Id)
	if err != nil {
		return fmt.Errorf("failed to get users in chat[%d], %s", chat.Id, err)
	}
	// inform if chat is open
	var errs []string
	for _, targetUser := range targetUsers {
		if targetUser == nil {
			continue
		}
		openChatId := state.GetOpenChat(targetUser.Id)
		if openChatId == 0 || openChatId != chat.Id {
			// TODO unread msg indicator -> counter update
			continue
		}
		uErr := distributeUpdateOfUser(state, targetUser, subjectUser, updateType)
		if uErr != nil {
			errs = append(errs, fmt.Sprintf("failed to distribute to user[%d], %s\n", targetUser.Id, uErr))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to distribute in chat[%d], %v", chat.Id, errs)
	}
	return nil
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
		return fmt.Errorf("subject user is nil")
	}
	log.Printf("userNameChanged TRACE informing target[%d] about subject[%d] name change\n", conn.User.Id, subject.Id)
	tmpl := subject.Template(0, 0, conn.User.Id)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template subject user")
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

package sse

import (
	"fmt"
	"log"

	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
)

// empty targetUser means all users in chat
func DistributeChat(
	state *state.State,
	db *d.DBConn,
	chat *app.Chat,
	author *app.User, // who made the change
	targetUser *app.User, // who to inform, nil for all users in chat
	subjectUser *app.User, // viewer, user affected by change, nil for every user in chat
	updateType event.EventType,
) error {
	if author == nil {
		return fmt.Errorf("author is nil")
	}
	if chat == nil {
		return fmt.Errorf("chat is nil")
	}
	if state == nil {
		return fmt.Errorf("state is nil")
	}

	var targetUsers []*app.User
	var err error
	if targetUser != nil {
		targetUsers = []*app.User{targetUser}
	} else {
		targetUsers, err = shared.GetChatUsers(db, chat.Id)
	}
	if err != nil {
		return fmt.Errorf("failed to get chat users: %s", err)
	}
	if len(targetUsers) <= 0 {
		return fmt.Errorf("chatUsers are empty in chat[%d]", chat.Id)
	}

	return distributeToUsers(state, chat, author, targetUsers, subjectUser, updateType)
}

func distributeToUsers(
	state *state.State,
	chat *app.Chat,
	author *app.User,
	targetUsers []*app.User,
	subjectUser *app.User,
	updateType event.EventType,
) error {
	var errors []string
	for _, targetUser := range targetUsers {
		log.Printf("distributeToUsers TRACE event[%v] about subject[%v] will be sent to user[%v] in chat[%v]\n",
			updateType, subjectUser, targetUser, chat)
		err := distributeChatToUser(
			state,
			author,
			targetUser,
			chat,
			subjectUser,
			updateType,
		)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("distributeToUsers errors: [%v]", errors)
	}
	return nil
}

func distributeChatToUser(
	state *state.State,
	author *app.User,
	targetUser *app.User,
	targetChat *app.Chat,
	subjectUser *app.User,
	updateType event.EventType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked: %v", r)
		}
	}()
	if author == nil {
		return fmt.Errorf("author is nil")
	}
	if targetChat == nil {
		return fmt.Errorf("targetChat is nil")
	}
	conns := state.GetConn(targetUser.Id)
	for _, conn := range conns {
		var connerr error

		switch updateType {
		case event.ChatAdd:
			connerr = chatCreate(conn, targetChat)
		case event.ChatInvite:
			connerr = chatInvite(conn, targetChat, author.Id, subjectUser)
		case event.ChatDrop:
			connerr = chatDelete(conn, targetChat.Id, targetChat.OwnerId, author.Id, targetUser.Id)
		case event.ChatExpel:
			connerr = chatExpel(conn, targetChat.Id, targetChat.OwnerId, author.Id, subjectUser.Id)
		case event.ChatLeave:
			connerr = chatLeave(conn, targetChat.Id, targetChat.OwnerId, author.Id, subjectUser.Id)
		case event.ChatClose:
			connerr = chatClose(conn, targetChat.Id, targetChat.OwnerId, author.Id, targetUser)
		default:
			connerr = fmt.Errorf("unknown event type[%v]", updateType)
		}

		if connerr != nil {
			log.Printf("distributeChatToUser ERROR failed to send update to user[%d], err[%s]\n",
				targetUser.Id, connerr)
			if err == nil {
				err = connerr
			} else {
				err = fmt.Errorf("%s, %s", err.Error(), connerr.Error())
			}
		}
	}
	return err
}

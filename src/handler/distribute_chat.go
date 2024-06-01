package handler

import (
	"fmt"
	"log"

	"prplchat/src/model/app"
	"prplchat/src/model/event"
)

// empty targetUser means all users in chat
func DistributeChat(
	state *AppState,
	chat *app.Chat,
	author *app.User, // who made the change
	targetUser *app.User, // who to inform, nil for all users in chat
	subjectUser *app.User, // which chat changed, nil for every user in chat
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

	var err error
	var targetUsers []*app.User
	if targetUser != nil {
		targetUsers = []*app.User{targetUser}
	} else {
		// have to get users by owner - author may have been removed
		targetUsers, err = chat.GetUsers(chat.Owner.Id)
	}

	if err != nil {
		return fmt.Errorf("fail to get users from chat[%d], %s", chat.Id, err)
	}
	if len(targetUsers) <= 0 {
		return fmt.Errorf("chatUsers are empty in chat[%d], %s", chat.Id, err)
	}

	return distributeToUsers(state, chat, author, targetUsers, subjectUser, updateType)
}

func distributeToUsers(
	state *AppState,
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
	state *AppState,
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
	var lmd FuncPerConn
	switch updateType {
	case event.ChatAdd:
		lmd = func(conn *Conn) error {
			return chatCreate(conn, targetChat, author)
		}
	case event.ChatInvite:
		lmd = func(conn *Conn) error {
			return chatInvite(conn, targetChat, author.Id, subjectUser)
		}
	case event.ChatDrop:
		lmd = func(conn *Conn) error {
			return chatDelete(conn, targetChat.Id, targetChat.Owner.Id, author.Id, targetUser.Id)
		}
	case event.ChatExpel:
		lmd = func(conn *Conn) error {
			return chatExpel(conn, targetChat.Id, targetChat.Owner.Id, author.Id, subjectUser.Id)
		}
	case event.ChatLeave:
		lmd = func(conn *Conn) error {
			return chatLeave(conn, targetChat.Id, targetChat.Owner.Id, author.Id, subjectUser.Id)
		}
	case event.ChatClose:
		lmd = func(conn *Conn) error {
			return chatClose(conn, targetChat.Id, targetChat.Owner.Id, author.Id, targetUser)
		}
	default:
		return fmt.Errorf("unknown event type[%v]", updateType)
	}
	conns := state.GetConn(targetUser.Id)
	for _, conn := range conns {
		connerr := lmd(conn)
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

package handler

import (
	"fmt"
	"log"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
)

// empty targetUser means all users in chat
func DistributeChat(
	state *AppState,
	chat *app.Chat,
	author *app.User, // who made the change
	targetUser *app.User, // who to inform, nil for all users in chat
	subjectUser *app.User, // which chat changed, nil for every user in chat
	evnt event.UpdateType,
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
		targetUsers, err = chat.GetUsers(author.Id)
	}

	if err != nil {
		return fmt.Errorf("get users, chat[%d], %s", chat.Id, err)
	}
	if len(targetUsers) == 0 {
		return fmt.Errorf("chatUsers are empty, chat[%d], %s", chat.Id, err)
	}

	return distributeToUsers(state, chat, author, targetUsers, subjectUser, evnt)
}

func distributeToUsers(
	state *AppState,
	chat *app.Chat,
	author *app.User,
	targetUsers []*app.User,
	subjectUser *app.User,
	evnt event.UpdateType,
) error {
	var errors []string
	for _, targetUser := range targetUsers {
		log.Printf("∞----> distributeToUsers TRACE event[%v] about subject[%v] will be sent to user[%v] in chat[%v]\n",
			evnt, subjectUser, targetUser, chat)
		err := distributeChatToUser(
			state,
			author,
			targetUser,
			chat,
			subjectUser,
			evnt,
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
	evnt event.UpdateType,
) error {
	if author == nil {
		return fmt.Errorf("author is nil")
	}
	if targetChat == nil {
		return fmt.Errorf("targetChat is nil")
	}

	conn, err := state.GetConn(targetUser.Id)
	if err != nil {
		return err
	}
	if conn.User.Id != targetUser.Id {
		return fmt.Errorf("user[%d] does not own conn[%v], user[%d] does", targetUser.Id, conn.Origin, conn.User.Id)
	}

	switch evnt {
	case event.ChatCreated:
		if author.Id != targetChat.Owner.Id {
			return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
				author.Id, targetChat.Owner.Id, targetChat.Id)
		}
		return chatCreate(conn, evnt, targetChat, author)

	case event.ChatInvite:
		if author.Id != targetChat.Owner.Id {
			return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
				author.Id, targetChat.Owner.Id, targetChat.Id)
		}
		if subjectUser == nil {
			return fmt.Errorf("subjectUser is nil for chatInvite")
		}
		return chatInvite(conn, evnt, targetChat, author.Id, subjectUser)

	case event.ChatDeleted:
		if author.Id != targetChat.Owner.Id && author.Id != targetUser.Id {
			return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
				author.Id, targetChat.Owner.Id, targetChat.Id)
		}
		if targetUser == nil {
			return fmt.Errorf("targetUser is nil for chatDeleted")
		}
		return chatDelete(conn, evnt, targetChat.Id, author.Id, targetUser.Id)

	case event.ChatExpel:
		if author.Id != targetChat.Owner.Id && author.Id != subjectUser.Id {
			return fmt.Errorf("author[%d] is not allowed to expel user[%d] from chat[%d]",
				author.Id, subjectUser.Id, targetChat.Id)
		}
		if subjectUser == nil {
			return fmt.Errorf("subjectUser is nil for chatExpel")
		}
		return chatExpel(conn, evnt, targetChat.Id, author.Id, subjectUser.Id)

	case event.ChatClose:
		if targetUser == nil {
			return fmt.Errorf("targetUser is nil for chatClose")
		}
		if author.Id != targetChat.Owner.Id && author.Id != targetUser.Id {
			return fmt.Errorf("author[%d] is not allowed to close chat[%d] for user[%d]",
				author.Id, targetChat.Id, targetUser.Id)
		}
		return chatClose(conn, evnt, targetChat.Id, author.Id, targetUser.Id)

	default:
		return fmt.Errorf("unknown event type[%v]", evnt)
	}
}

package pub

import (
	"fmt"
	"log"
	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
	"neon-chat/src/state"

	"github.com/jmoiron/sqlx"
)

func InviteUser(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	inviteeName string,
) (*app.Chat, *app.User, error) {
	appInvitee, err := searchUser(dbConn.Tx, inviteeName)
	if err != nil {
		log.Printf("ERROR HandleUserInvite invitee not found [%s], %s\n", inviteeName, err.Error())
		return nil, nil, fmt.Errorf("invitee not found")
	} else if appInvitee == nil {
		log.Printf("WARN HandleUserInvite invitee not found [%s]\n", inviteeName)
		return nil, nil, nil
	}
	appChat, err := priv.GetChat(state, dbConn.Tx, user, chatId)
	if err != nil {
		log.Printf("ERROR HandleUserInvite user[%d] cannot invite into chat[%d], %s\n",
			user.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("cannot find chat: %s", err.Error())
	} else if appChat == nil {
		log.Printf("WARN HandleUserInvite user[%d] cannot invite into chat[%d]\n", user.Id, chatId)
		return nil, nil, fmt.Errorf("chat not found")
	}
	err = db.AddChatUser(dbConn.Tx, chatId, appInvitee.Id)
	if err != nil {
		log.Printf("ERROR HandleUserInvite failed to add user[%d] to chat[%d] in db, %s\n",
			appInvitee.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to add user to chat in db")
	}
	return appChat, appInvitee, nil
}

func searchUser(dbConn sqlx.Ext, userName string) (*app.User, error) {
	log.Printf("TRACE FindUser IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(dbConn, userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}
	var dbAvatar *db.Avatar
	if dbUser != nil {
		dbAvatar, _ = db.GetAvatar(dbConn, dbUser.Id)
	}

	log.Printf("TRACE FindUser OUT user[%s]\n", userName)
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}

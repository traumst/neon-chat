package pub

import (
	"fmt"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func InviteUser(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	inviteeName string,
) (*app.Chat, *app.User, error) {
	appInvitee, err := priv.SearchUser(dbConn.Tx, inviteeName)
	if err != nil {
		log.Printf("HandleUserInvite ERROR invitee not found [%s], %s\n", inviteeName, err.Error())
		return nil, nil, fmt.Errorf("invitee not found")
	} else if appInvitee == nil {
		log.Printf("HandleUserInvite WARN invitee not found [%s]\n", inviteeName)
		return nil, nil, nil
	}
	appChat, err := priv.GetChat(state, dbConn.Tx, user, chatId)
	if err != nil {
		log.Printf("HandleUserInvite ERROR user[%d] cannot invite into chat[%d], %s\n",
			user.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("cannot find chat: %s", err.Error())
	} else if appChat == nil {
		log.Printf("HandleUserInvite WARN user[%d] cannot invite into chat[%d]\n", user.Id, chatId)
		return nil, nil, fmt.Errorf("chat not found")
	}
	err = db.AddChatUser(dbConn.Tx, chatId, appInvitee.Id)
	if err != nil {
		log.Printf("HandleUserInvite ERROR failed to add user[%d] to chat[%d] in db, %s\n",
			appInvitee.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to add user to chat in db")
	}
	return appChat, appInvitee, nil
}

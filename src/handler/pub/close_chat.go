package pub

import (
	"fmt"

	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func CloseChat(state *state.State, dbConn *db.DBConn, user *app.User, chatId uint) (string, error) {
	err := state.CloseChat(user.Id, chatId)
	if err != nil {
		return "", fmt.Errorf("close chat[%d] for user[%d]: %s", chatId, user.Id, err)
	}
	return TemplateWelcome(user)
}

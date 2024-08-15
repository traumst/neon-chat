package shared

import (
	"log"
	"net/http"
	"prplchat/src/convert"
	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func TemplateHome(
	state *state.State,
	db *db.DBConn,
	r *http.Request,
	user *app.User,
) (string, error) {
	var avatarTmpl *template.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	openChatTemplate := TemplateOpenChat(state, db, user)
	chats, err := handler.GetChats(state, db, user.Id)
	if err != nil {
		log.Printf("[%s] templateHome ERROR, failed getting chats for user[%d], %s\n",
			h.GetReqId(r), user.Id, err.Error())
		return "", err
	}
	chatTemplates := make([]*template.ChatTemplate, 0)
	for _, chat := range chats {
		chatTemplates = append(chatTemplates, chat.Template(user, user, []*app.User{}, []*app.Message{}))
	}
	var openChatId uint
	var chatOwnerId uint
	if openChatTemplate != nil {
		openChatId = openChatTemplate.ChatId
		chatOwnerId = openChatTemplate.Owner.UserId
	}
	userTemplate := user.Template(openChatId, chatOwnerId, user.Id)
	home := template.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      openChatTemplate,
		User:          *userTemplate,
		IsAuthorized:  true,
		LoginTemplate: template.AuthTemplate{},
		Avatar:        avatarTmpl,
	}
	return home.HTML()
}

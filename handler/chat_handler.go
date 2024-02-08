package handler

import (
	"fmt"
	"log"
	"net/http"

	"go.chat/model"
	"go.chat/utils"
)

func SendChat(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	if up.Msg != nil {
		log.Printf("--%s-- sendChat INFO unexpected input msg %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- sendChat TRACE templating %s\n", utils.GetReqId(r), up.Chat.Log())
	html, err := up.Chat.ToTemplate(user).GetShortHTML()
	if err != nil {
		log.Printf("--%s-- sendChat ERROR templating %s, %s\n",
			utils.GetReqId(r), up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.MessageEventName, utils.RandStringBytes(5))
	distribute(w, r, up.Chat, model.ChatEventName, eventID, user, html)
	log.Printf("--%s-- sendChat TRACE message from %s\n", utils.GetReqId(r), user)
}

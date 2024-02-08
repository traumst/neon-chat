package handler

import (
	"fmt"
	"log"
	"net/http"

	"go.chat/model"
	"go.chat/utils"
)

func SendMessage(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	if up.Msg == nil {
		log.Printf("--%s-- sendMessage INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- sendMessage TRACE templating %s\n", utils.GetReqId(r), up.Msg.Log())
	html, err := up.Msg.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("--%s-- sendMessage ERROR templating %s, %s\n",
			utils.GetReqId(r), up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.MessageEventName, utils.RandStringBytes(5))
	distribute(w, r, up.Chat, model.MessageEventName, eventID, user, html)
	log.Printf("--%s-- sendMessage TRACE message from %s\n", utils.GetReqId(r), user)
}

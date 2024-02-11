package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.chat/model"
	"go.chat/utils"
)

func SendMessage(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	if up.Msg == nil {
		log.Printf("--%s-- SendMessage INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- SendMessage TRACE templating %s\n", utils.GetReqId(r), up.Msg.Log())
	html, err := up.Msg.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("--%s-- SendMessage ERROR templating %s, %s\n",
			utils.GetReqId(r), up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.MessageEventName, utils.RandStringBytes(5))
	writeMessage(w, r, up.Chat, model.MessageEventName, eventID, user, html)
	log.Printf("--%s-- SendMessage TRACE message from %s\n", utils.GetReqId(r), user)
}

func writeMessage(
	w *http.ResponseWriter,
	r *http.Request,
	chat *model.Chat,
	eventName model.SSEvent,
	eventID string,
	user string,
	html string) {
	log.Printf("---- writeMessage TRACE IN type[%s], event[%s] by user [%s]\n", eventName, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")
	html = strings.ReplaceAll(html, "  ", " ")

	log.Printf("--%s-- writeMessage TRACE html to %s\n", utils.GetReqId(r), user)
	fmt.Fprintf(*w, "id: %s\n\n", eventID)
	fmt.Fprintf(*w, "event: %s\n", eventName)
	fmt.Fprintf(*w, "data: %s\n\n", html)
	(*w).(http.Flusher).Flush()

	log.Printf("---- writeMessage TRACE OUT msg[%s] event[%s] by user [%s]\n", eventName, eventID, user)
}

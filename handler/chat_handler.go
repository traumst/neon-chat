package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.chat/model"
	"go.chat/utils"
)

func SendChat(
	w *http.ResponseWriter,
	up model.UserUpdate,
	user string) {
	if up.Msg != nil {
		log.Printf("∞----∞ SendChat INFO unexpected input msg %s\n", up.Msg.Log())
		return
	}

	log.Printf("∞----∞ SendChat TRACE templating %s\n", up.Chat.Log())
	html, err := up.Chat.ToTemplate(user).GetShortHTML()
	if err != nil {
		log.Printf("∞----∞ SendChat ERROR templating %s, %s\n", up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.ChatEventName, utils.RandStringBytes(5))
	writeChat(w, up.Chat, model.ChatEventName, eventID, user, html)
}

func writeChat(
	w *http.ResponseWriter,
	chat *model.Chat,
	eventName model.SSEvent,
	eventID string,
	user string,
	html string) {
	log.Printf("∞----∞ writeChat TRACE IN type[%s], event[%s] by user [%s]\n", eventName, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")

	if w == nil {
		log.Printf("∞----∞ writeChat ERROR writer is nil\n")
		return
	}

	log.Printf("∞----∞ writeChat TRACE html to %s\n", user)
	fmt.Fprintf(*w, "id: %s\n\n", eventID)
	fmt.Fprintf(*w, "event: %s\n", eventName)
	fmt.Fprintf(*w, "data: %s\n\n", html)
	(*w).(http.Flusher).Flush()

	log.Printf("∞----∞ writeChat TRACE OUT msg[%s] event[%s] by user [%s]\n", eventName, eventID, user)
}

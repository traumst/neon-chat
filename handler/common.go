package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.chat/model"
	"go.chat/utils"
)

func distribute(
	w *http.ResponseWriter,
	r *http.Request,
	chat *model.Chat,
	msgType model.SSEvent,
	eventID string,
	user string,
	html string) {
	log.Printf("---- distribute TRACE IN type[%s], event[%s] by user [%s]\n", msgType, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")
	html = strings.ReplaceAll(html, "  ", " ")

	chatUsers, err := chat.GetUsers(user)
	if err != nil || chatUsers == nil {
		log.Printf("---- distribute ERROR get users, %s\n", err)
		return
	}

	for _, u := range chatUsers {
		if u == user && msgType == model.MessageEventName {
			continue
		}

		flush(w, r, msgType, user, eventID, html)
	}
	log.Printf("---- distribute TRACE OUT msg[%s], event[%s] by user [%s]\n", msgType, eventID, user)
}

func flush(
	w *http.ResponseWriter,
	r *http.Request,
	eventName model.SSEvent,
	user string,
	eventID string,
	html string) {
	log.Printf("--%s-- flush TRACE html to %s\n", utils.GetReqId(r), user)
	fmt.Fprintf(*w, "id: %s\n\n", eventID)
	fmt.Fprintf(*w, "event: %s\n", eventName)
	fmt.Fprintf(*w, "data: %s\n\n", html)
	(*w).(http.Flusher).Flush()
}

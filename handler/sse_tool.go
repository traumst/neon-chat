package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.chat/model"
	"go.chat/utils"
)

func SendEvent(
	reqId string,
	w *http.ResponseWriter,
	up model.UserUpdate,
	user string) {
	log.Printf("∞--%s--> SendEvent TRACE IN %+v\n", reqId, up)
	if up.User == "" || up.Msg == "" {
		log.Printf("<--%s--∞ SendEvent ERROR user or msg is empty, %+v\n", reqId, up)
		return
	}

	var eventName model.SSEvent
	switch up.Type {
	case model.ChatUpdate, model.ChatInvite:
		eventName = model.ChatEventName
	case model.MessageUpdate:
		eventName = model.MessageEventName
	case model.PingUpdate:
		eventName = model.PingEventName
	default:
		eventName = model.Unknown
		log.Printf("<--%s--∞ SendEvent ERROR unknown event type, %+v\n", reqId, up)
		return
	}

	eventID := fmt.Sprintf("%s-%s", eventName, utils.RandStringBytes(5))
	// must escape newlines in SSE
	html := strings.ReplaceAll(up.Msg, "\n", " ")
	log.Printf("<--%s--∞ SendEvent TRACE type[%s], event[%s], html[%s]\n", reqId, eventName, eventID, html)
	if w == nil {
		log.Printf("<--%s--∞ SendEvent ERROR writer is nil on event[%s]\n", reqId, eventName)
		return
	}

	sendSse(reqId, w, eventID, string(eventName), html)
	log.Printf("<--%s--∞ SendEvent TRACE OUT %+v\n", reqId, up)
}

func sendSse(
	reqId string,
	w *http.ResponseWriter,
	eventID string,
	eventName string,
	html string) {
	log.Printf("<---%s---∞ sendSse TRACE id[%s], name[%s]\n", reqId, eventID, eventName)
	fmt.Fprintf(*w, "id: %s\n\n", eventID)
	fmt.Fprintf(*w, "event: %s\n", eventName)
	fmt.Fprintf(*w, "data: %s\n\n", html)
	(*w).(http.Flusher).Flush()
}

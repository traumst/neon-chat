package model

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.chat/utils"
)

type App struct {
	State AppState
}

func (app *App) PollUpdatesForUser(conn *Conn, user string) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%s]\n", conn.Origin, conn.User)
	channel := conn.Channel
	origin := conn.Origin
	reader := conn.Reader
	for {
		log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%s] is waiting for updates\n", origin, conn.User)
		select {
		case <-reader.Context().Done():
			return
		case update := <-channel:
			if update.Author == "" || update.Msg == "" {
				log.Printf("<--%s--∞ APP.PollUpdatesForUser INFO user or msg is empty, update[%+v]\n", origin, update)
				return
			}
			app.sendUpdates(conn, update, user)
		}
	}
}

func (app *App) sendUpdates(conn *Conn, up UserUpdate, user string) {
	log.Printf("∞--%s--> APP.sendUpdates TRACE IN [%s], input[%+v]\n", conn.Origin, user, up)
	chat, err := app.State.GetChat(user, up.ChatID)
	if err != nil {
		log.Printf("<--%s--∞ APP.sendUpdates WARN %s\n", conn.Origin, err)
		return
	}
	users, err := chat.GetUsers(user)
	if err != nil {
		log.Printf("<--%s--∞ APP.sendUpdates ERROR chat[%s], %s\n",
			conn.Origin, chat.Name, err)
		return
	}
	for _, u := range users {
		isSent := trySend(conn.Origin, conn, up, u)
		if !isSent {
			log.Printf("<--%s--∞ APP.sendUpdates ERROR failed to send update to user[%s]\n", conn.Origin, u)
		}
	}
	log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%s]\n", conn.Origin, user)
}

func trySend(reqId string, conn *Conn, up UserUpdate, user string) bool {
	w := conn.Writer

	if user == "" || up.Msg == "" {
		log.Printf("<--%s--∞ trySend ERROR user or msg is empty, user[%s], msg[%s]\n", reqId, user, up.Msg)
		return false
	}
	if w == nil {
		log.Printf("<--%s--∞ trySend ERROR writer is nil\n", reqId)
		return false
	}

	switch up.Type {
	case ChatUpdate, ChatInvite, MessageUpdate:
		log.Printf("∞--%s--> trySend TRACE sending event[%s] to user[%s] via w[%T]\n", reqId, up.Type, user, w)
		err := sendEvent(&w, up.Type.String(), up.Msg)
		if err != nil {
			log.Printf("<--%s--∞ trySend ERROR failed to send event[%s] to user[%s], %s\n", reqId, up.Type, user, err)
			return false
		}
	default:
		log.Printf("<--%s--∞ trySend ERROR unknown update event[%s], update[%+v]\n", reqId, up.Type, up)
		return false
	}

	log.Printf("<--%s--∞ trySend TRACE event[%s] sent to user[%s]\n", reqId, up.Type, user)
	return true
}

func sendEvent(w *http.ResponseWriter, eventName string, html string) error {
	writer := *w
	eventID := utils.RandStringBytes(5)
	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")

	_, err := fmt.Fprintf(writer, "id: %s\n\n", eventID)
	if err != nil {
		return fmt.Errorf("failed to write id[%s]", eventID)
	}
	_, err = fmt.Fprintf(writer, "event: %s\n", eventName)
	if err != nil {
		return fmt.Errorf("failed to write event[%s]", eventName)
	}
	_, err = fmt.Fprintf(writer, "data: %s\n\n", html)
	if err != nil {
		return fmt.Errorf("failed to write data[%s]", html)
	}
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		return fmt.Errorf("writer does not support flushing")
	}
	flusher.Flush()
	return nil
}

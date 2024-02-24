package model

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"go.chat/utils"
)

type App struct {
	State       AppState
	pollCounter atomic.Int32
}

func (app *App) PollUpdatesForUser(conn *Conn, pollingUser string) {
	app.pollCounter.Add(1)
	defer app.pollCounter.Add(-1)
	origin := conn.Origin
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE pollCounter[%d]\n", origin, app.pollCounter.Load())
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%s]\n", origin, conn.User)
	for {
		log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%s] is waiting for updates\n", origin, conn.User)
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%s] conn[%v] disonnected\n",
				origin, pollingUser, conn.Origin)
			return
		case update := <-conn.In:
			//defer func() { conn.Out <- update }()
			if conn.User != pollingUser {
				log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%v] is does not own conn[%v]\n",
					origin, pollingUser, conn)
				conn.Out <- update
				continue
			}
			if update.Author == "" || update.RawHtml == "" {
				log.Printf("<--%s--∞ APP.PollUpdatesForUser INFO user or msg is empty, update[%+v]\n", origin, update)
				conn.Out <- update
				continue
			}
			app.sendUpdates(conn, update, pollingUser)
		}
	}
}

func (app *App) sendUpdates(conn *Conn, up UserUpdate, pollingUser string) {
	log.Printf("∞--%s--> APP.sendUpdates TRACE IN [%s], input[%+v]\n", conn.Origin, pollingUser, up)

	isSent := trySend(conn.Origin, conn, up, pollingUser)
	if isSent {
		up.RawHtml = fmt.Sprintf("SENT TO: %s", pollingUser)
		conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%s]\n", conn.Origin, pollingUser)
	} else {
		up.RawHtml = fmt.Sprintf("ERROR SENDING TO: %s", pollingUser)
		conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates ERROR failed to send update to user[%s]\n", conn.Origin, pollingUser)
	}
}

func trySend(reqId string, conn *Conn, up UserUpdate, user string) bool {
	w := conn.Writer

	if user == "" || up.RawHtml == "" {
		log.Printf("<--%s--∞ trySend ERROR user or msg is empty, user[%s], msg[%s]\n", reqId, user, up.RawHtml)
		return false
	}
	if w == nil {
		log.Printf("<--%s--∞ trySend ERROR writer is nil\n", reqId)
		return false
	}

	switch up.Type {
	case ChatUpdate, ChatInvite, MessageUpdate:
		log.Printf("∞--%s--> trySend TRACE sending event[%s] to user[%s] via w[%T]\n", reqId, up.Type, user, w)
		err := sendEvent(&w, up.Type.String(), up.RawHtml)
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

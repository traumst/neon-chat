package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"go.chat/model"
	"go.chat/utils"
)

func PollUpdatesForUser(conn *model.Conn, pollingUser string) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%s]\n", conn.Origin, conn.User)
	var wg sync.WaitGroup
	for {
		log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%s] is waiting for updates\n", conn.Origin, conn.User)
		done := false
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%s] conn[%v] disonnected\n",
				conn.Origin, pollingUser, conn.Origin)
			done = true
		case up := <-conn.In:
			wg.Add(1)
			go func() {
				defer wg.Done()
				sendUpdates(conn, up, pollingUser)
			}()
		}

		if done {
			break
		}
	}
	wg.Wait()
}

func sendUpdates(conn *model.Conn, up model.LiveUpdate, pollingUser string) {
	log.Printf("∞--%s--> APP.sendUpdates TRACE IN [%s], input[%+v]\n", conn.Origin, pollingUser, up)
	origin := conn.Origin
	if conn.User != pollingUser {
		log.Printf("<--%s--∞ APP.sendUpdates WARN user[%v] is does not own conn[%v]\n", origin, pollingUser, conn)
		return
	}
	if up.Author == "" || up.Data == "" {
		log.Printf("<--%s--∞ APP.sendUpdates INFO user or msg is empty, update[%+v]\n", origin, up)
		return
	}
	isSent := trySend(origin, conn, up, pollingUser)
	if isSent {
		up.Error = fmt.Errorf("SENT TO: %s", pollingUser)
		conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%s]\n", origin, pollingUser)
	} else {
		up.Error = fmt.Errorf("ERROR SENDING TO: %s", pollingUser)
		conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates ERROR failed to send update to user[%s]\n", origin, pollingUser)
	}
}

func trySend(reqId string, conn *model.Conn, up model.LiveUpdate, user string) bool {
	w := conn.Writer
	if user == "" || up.Data == "" {
		log.Printf("<--%s--∞ trySend ERROR user or msg is empty, user[%s], msg[%s]\n", reqId, user, up.Data)
		return false
	}
	if w == nil {
		log.Printf("<--%s--∞ trySend ERROR writer is nil\n", reqId)
		return false
	}
	event := up.Event.String()
	switch up.Event {
	case model.ChatCreated,
		model.ChatInvite,
		model.ChatDeleted,
		model.MessageAdded,
		model.MessageDeleted:
		log.Printf("∞--%s--> trySend TRACE sending event[%s] to user[%s] via w[%T]\n", reqId, event, user, w)
		err := sendEvent(&w, event, up.Data)
		if err != nil {
			log.Printf("<--%s--∞ trySend ERROR failed to send event[%s] to user[%s], %s\n", reqId, event, user, err)
			return false
		}
	default:
		log.Printf("<--%s--∞ trySend ERROR unknown update event[%s], update[%+v]\n", reqId, event, up)
		return false
	}
	log.Printf("<--%s--∞ trySend TRACE event[%s] sent to user[%s]\n", reqId, event, user)
	return true
}

func sendEvent(w *http.ResponseWriter, eventName string, html string) error {
	writer := *w
	eventID := utils.RandStringBytes(5)
	_, err := fmt.Fprintf(writer, "id: %s\n\n", eventID)
	if err != nil {
		return fmt.Errorf("failed to write id[%s]", eventID)
	}
	_, err = fmt.Fprintf(writer, "event: %s\n", eventName)
	if err != nil {
		return fmt.Errorf("failed to write event[%s]", eventName)
	}
	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")
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

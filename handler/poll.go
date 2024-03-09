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
	log.Printf("∞--%s--> APP.sendUpdates TRACE IN [%s], input[%s]\n", conn.Origin, pollingUser, up.String())
	origin := conn.Origin
	if conn.User != pollingUser {
		log.Printf("<--%s--∞ APP.sendUpdates WARN user[%v] is does not own conn[%v]\n", origin, pollingUser, conn)
		return
	}
	if up.Author == "" || up.Data == "" {
		log.Printf("<--%s--∞ APP.sendUpdates INFO user or msg is empty, update[%v]\n", origin, up)
		return
	}
	err := trySend(conn, up, pollingUser)
	if err != nil {
		up.Error = fmt.Errorf("ERROR SENDING TO: %s", pollingUser)
		conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates ERROR failed to send update to user[%s]\n", origin, pollingUser)
		return
	}
	up.Error = fmt.Errorf("SENT TO: %s", pollingUser)
	conn.Out <- up
	log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%s]\n", origin, pollingUser)
}

func trySend(conn *model.Conn, up model.LiveUpdate, user string) error {
	w := conn.Writer
	if user == "" || up.Data == "" {
		return fmt.Errorf("trySend ERROR user or msg is empty, user[%s], msg[%s]", user, up.Data)
	}
	if w == nil {
		return fmt.Errorf("trySend ERROR writer is nil")
	}
	event := up.Event.String()
	switch up.Event {
	case model.ChatCreated, model.ChatInvite, model.MessageAdded:
		err := sendEvent(&w, event, up.Data)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to send event[%s] to user[%s], %s", event, user, err)
		}
	case model.MessageDeleted:
		err := deleteMsg(&w, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete message to user[%s], %s", user, err)
		}
	case model.ChatDeleted:
		err := deleteChat(&w, up.Data)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete chat to user[%s], %s", user, err)
		}
	default:
		return fmt.Errorf("trySend ERROR unknown update event[%s], update[%+v]", event, up)
	}
	return nil
}

func deleteMsg(w *http.ResponseWriter, up model.LiveUpdate) error {
	dropEvent := fmt.Sprintf("%s-chat-%d-msg-%d", model.MessageDropEventName, up.ChatID, up.MsgID)
	log.Printf("deleteMsg TRACE in, sse[%s]\n", dropEvent)
	err := sendEvent(w, dropEvent, "[deleted]")
	if err != nil {
		return fmt.Errorf("deleteMsg ERROR failed to send message-drop event, %s", err)
	}
	return nil
}

func deleteChat(w *http.ResponseWriter, data string) error {
	log.Printf("deleteChat TRACE in\n")
	err := sendEvent(w, string(model.ChatCloseEventName), data)
	if err != nil {
		return fmt.Errorf("deleteChat ERROR failed to send chat-close event, %s", err)
	}
	err = sendEvent(w, string(model.ChatDropEventName), "")
	if err != nil {
		return fmt.Errorf("deleteChat ERROR failed to send chat-drop event, %s", err)
	}
	return nil
}

func sendEvent(w *http.ResponseWriter, eventName string, html string) error {
	writer := *w
	eventID := utils.RandStringBytes(5)
	_, err := fmt.Fprintf(writer, "id: %s\n", eventID)
	if err != nil {
		return fmt.Errorf("failed to write id[%s]", eventID)
	}
	_, err = fmt.Fprintf(writer, "event: %s\n", eventName)
	if err != nil {
		return fmt.Errorf("failed to write event[%s]", eventName)
	}
	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")
	// remove double spaces
	for strings.Contains(html, "  ") {
		html = strings.ReplaceAll(html, "  ", " ")
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

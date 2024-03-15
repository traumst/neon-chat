package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"go.chat/model"
	"go.chat/model/event"
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
				conn.Out <- up
			}()
		}

		if done {
			break
		}
	}
	wg.Wait()
}

func sendUpdates(conn *model.Conn, up event.LiveUpdate, pollingUser string) {
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
	log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%s]\n", origin, pollingUser)
}

func trySend(conn *model.Conn, up event.LiveUpdate, user string) error {
	w := conn.Writer
	if user == "" || up.Data == "" {
		return fmt.Errorf("trySend ERROR user or msg is empty, user[%s], msg[%s]", user, up.Data)
	}
	if w == nil {
		return fmt.Errorf("trySend ERROR writer is nil")
	}
	switch up.Event {
	case event.ChatCreated:
		err := SSEvent(&w, event.ChatAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to send to user[%s], %s", user, err)
		}
	case event.ChatInvite:
		if up.Author == up.UserID {
			return nil
		}
		err := SSEvent(&w, event.ChatAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to send to user[%s], %s", user, err)
		}
	case event.ChatUserDrop:
		if up.Author == up.UserID {
			return fmt.Errorf("trySend ERROR user[%s] is trying to drop itself from chat[%d]", user, up.ChatID)
		}
		err := SSEvent(&w, event.ChatUserDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to drop user from chat to user[%s], %s", user, err)
		}
	case event.MessageAdded:
		err := SSEvent(&w, event.MessageAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to add message to user[%s], %s", user, err)
		}
	case event.MessageDeleted:
		err := SSEvent(&w, event.MessageDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete message to user[%s], %s", user, err)
		}
	case event.ChatDeleted:
		err := SSEvent(&w, event.ChatDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete chat to user[%s], %s", user, err)
		}
	case event.ChatClose:
		err := SSEvent(&w, event.ChatCloseEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to close chat to user[%s], %s", user, err)
		}
	default:
		return fmt.Errorf("trySend ERROR unknown update event[%v], update[%s]", up.Event, up.String())
	}
	return nil
}

func SSEvent(w *http.ResponseWriter, event event.SSEvent, up event.LiveUpdate) error {
	if up.ChatID < 0 {
		panic("ChatID should not be empty")
	}
	if up.UserID == "" {
		panic("UserID should not be empty")
	}
	eventName := Format(event, up.ChatID, up.UserID, up.MsgID)
	eventID := utils.RandStringBytes(5)
	data := trim(up.Data)
	writer := *w
	_, err := fmt.Fprintf(writer, "id: %s\n", eventID)
	if err != nil {
		return fmt.Errorf("failed to write id[%s]", eventID)
	}
	_, err = fmt.Fprintf(writer, "event: %s\n", eventName)
	if err != nil {
		return fmt.Errorf("failed to write event[%s]", eventName)
	}
	_, err = fmt.Fprintf(writer, "data: %s\n\n", data)
	if err != nil {
		return fmt.Errorf("failed to write data[%s]", data)
	}
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		return fmt.Errorf("writer does not support flushing")
	}
	flusher.Flush()
	return nil
}

func Format(e event.SSEvent, chatID int, userID string, msgID int) string {
	switch e {
	case event.ChatAddEventName:
		return string(event.ChatAddEventName)
	case event.ChatDropEventName:
		return fmt.Sprintf("%s-%d", event.ChatDropEventName, chatID)
	case event.ChatCloseEventName:
		return fmt.Sprintf("%s-%d", event.ChatCloseEventName, chatID)
	case event.ChatUserDropEventName:
		return fmt.Sprintf("%s-%d-user-%s", event.ChatUserDropEventName, chatID, userID)
	case event.MessageAddEventName:
		return fmt.Sprintf("%s-chat-%d", event.MessageAddEventName, chatID)
	case event.MessageDropEventName:
		return fmt.Sprintf("%s-chat-%d-msg-%d", event.MessageDropEventName, chatID, msgID)
	default:
		panic(fmt.Sprintf("unknown event type[%v]", e))
	}
}

func trim(s string) string {
	// must escape newlines in SSE
	res := strings.ReplaceAll(s, "\n", " ")
	// remove double spaces
	for strings.Contains(res, "  ") {
		res = strings.ReplaceAll(res, "  ", " ")
	}
	return res
}

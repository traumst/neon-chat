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

func PollUpdatesForUser(conn *model.Conn, pollingUserId uint) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%d]\n", conn.Origin, conn.User.Id)
	var wg sync.WaitGroup
	done := false
	for !done {
		log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%d] is waiting for updates\n", conn.Origin, conn.User.Id)
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			wg.Add(1)
			go func() {
				defer wg.Done()
				sendUpdates(conn, up, pollingUserId)
				//conn.Out <- up
			}()
		}
	}
	wg.Wait()
}

func sendUpdates(conn *model.Conn, up event.LiveUpdate, pollingUserId uint) {
	log.Printf("∞--%s--> APP.sendUpdates TRACE IN user[%d], input[%s]\n", conn.Origin, pollingUserId, up.String())
	origin := conn.Origin
	if conn.User.Id != pollingUserId {
		log.Printf("<--%s--∞ APP.sendUpdates WARN user[%v] is does not own conn[%v]\n", origin, pollingUserId, conn)
		return
	}
	err := trySend(conn, up)
	if err != nil {
		up.Error = fmt.Errorf("ERROR SENDING TO: %d", pollingUserId)
		//conn.Out <- up
		log.Printf("<--%s--∞ APP.sendUpdates ERROR failed to send update to user[%d], err[%s]\n",
			origin, pollingUserId, err)
		return
	}
	log.Printf("<--%s--∞ APP.sendUpdates TRACE OUT user[%d]\n", origin, pollingUserId)
}

func trySend(conn *model.Conn, up event.LiveUpdate) error {
	w := conn.Writer
	if up.UserId == 0 {
		return fmt.Errorf("trySend ERROR user is empty, user[%d], msg[%s]", up.UserId, up.Data)
	}
	if w == nil {
		return fmt.Errorf("trySend ERROR writer is nil")
	}
	switch up.Event {
	case event.MessageAdded:
		err := SSEvent(&w, event.MessageAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to add message to user[%d], %s", up.UserId, err)
		}
	case event.MessageDeleted:
		err := SSEvent(&w, event.MessageDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete message to user[%d], %s", up.UserId, err)
		}
	case event.ChatCreated, event.ChatInvite:
		err := SSEvent(&w, event.ChatAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to send to user[%d], %s", up.UserId, err)
		}
	case event.ChatExpel:
		err := SSEvent(&w, event.ChatExpelEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to drop user from chat to user[%d], %s", up.UserId, err)
		}
	case event.ChatClose:
		err := SSEvent(&w, event.ChatCloseEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to close chat to user[%d], %s", up.UserId, err)
		}
	case event.ChatDeleted:
		err := SSEvent(&w, event.ChatDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete chat to user[%d], %s", up.UserId, err)
		}
	default:
		return fmt.Errorf("trySend ERROR unknown update event[%v], update[%s]", up.Event, up.String())
	}
	return nil
}

func SSEvent(w *http.ResponseWriter, event event.SSEvent, up event.LiveUpdate) error {
	if up.ChatId < 0 {
		panic("ChatID should not be empty")
	}
	if up.UserId == 0 {
		panic("UserID should not be empty")
	}
	eventName := Format(event, up.ChatId, up.UserId, up.MsgId)
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

func Format(e event.SSEvent, chatId int, userId uint, msgId int) string {
	switch e {
	case event.MessageAddEventName:
		return fmt.Sprintf("%s-chat-%d", event.MessageAddEventName, chatId)
	case event.MessageDropEventName:
		return fmt.Sprintf("%s-chat-%d-msg-%d", event.MessageDropEventName, chatId, msgId)
	case event.ChatAddEventName:
		return string(event.ChatAddEventName)
	case event.ChatExpelEventName:
		return fmt.Sprintf("%s-%d-user-%d", event.ChatExpelEventName, chatId, userId)
	case event.ChatDropEventName:
		return fmt.Sprintf("%s-%d", event.ChatDropEventName, chatId)
	case event.ChatCloseEventName:
		return fmt.Sprintf("%s-%d", event.ChatCloseEventName, chatId)
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

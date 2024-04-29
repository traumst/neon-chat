package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/src/model/event"
	"go.chat/src/utils"
)

func PollUpdatesForUser(conn *Conn, pollingUserId uint) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%d]\n", conn.Origin, conn.User.Id)
	var wg sync.WaitGroup
	done := false
	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%d] is receiving update[%s]\n",
					conn.Origin, conn.User.Id, up.Event)
				sendUpdates(conn, up, pollingUserId)
				//conn.Out <- up
			}()
		}
	}
	wg.Wait()
}

func sendUpdates(conn *Conn, up event.LiveUpdate, pollingUserId uint) {
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

func trySend(conn *Conn, up event.LiveUpdate) error {
	w := conn.Writer
	if up.UserId == 0 {
		return fmt.Errorf("trySend ERROR user is empty, user[%d], msg[%s]", up.UserId, up.Data)
	}
	if w == nil {
		return fmt.Errorf("trySend ERROR writer is nil")
	}
	switch up.Event {
	case event.MessageAdded:
		err := flushEvent(&w, event.MessageAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to add message to user[%d], %s", up.UserId, err)
		}
	case event.MessageDeleted:
		err := flushEvent(&w, event.MessageDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete message to user[%d], %s", up.UserId, err)
		}
	case event.ChatCreated, event.ChatInvite:
		err := flushEvent(&w, event.ChatAddEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to send to user[%d], %s", up.UserId, err)
		}
	case event.ChatExpel:
		err := flushEvent(&w, event.ChatExpelEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to drop user from chat to user[%d], %s", up.UserId, err)
		}
	case event.ChatClose:
		err := flushEvent(&w, event.ChatCloseEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to close chat to user[%d], %s", up.UserId, err)
		}
	case event.ChatDeleted:
		err := flushEvent(&w, event.ChatDropEventName, up)
		if err != nil {
			return fmt.Errorf("trySend ERROR failed to delete chat to user[%d], %s", up.UserId, err)
		}
	default:
		return fmt.Errorf("trySend ERROR unknown update event[%v], update[%s]", up.Event, up.String())
	}
	return nil
}

func flushEvent(w *http.ResponseWriter, evnt event.SSEvent, up event.LiveUpdate) error {
	if up.ChatId < 0 {
		panic("ChatId should not be empty")
	}
	if up.UserId == 0 {
		panic("UserId should not be empty")
	}
	eventName := evnt.Format(up.ChatId, up.UserId, up.MsgId)
	eventId := utils.RandStringBytes(5)
	// must escape newlines in SSE
	data := utils.Trim(up.Data)
	_, err := fmt.Fprintf(*w, "id: %s\n", eventId)
	if err != nil {
		return fmt.Errorf("failed to write id[%s]", eventId)
	}
	_, err = fmt.Fprintf(*w, "event: %s\n", eventName)
	if err != nil {
		return fmt.Errorf("failed to write event[%s]", eventName)
	}
	_, err = fmt.Fprintf(*w, "data: %s\n\n", data)
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

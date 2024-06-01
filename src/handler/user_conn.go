package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

type Conn struct {
	Id     int
	User   *app.User
	Origin string
	Writer http.ResponseWriter
	Reader http.Request
	In     chan event.LiveEvent // TODO test load
	//Out    chan event.LiveUpdate
}

// TODO should be LRUCache
type ActiveConnections map[uint][]*Conn

var mu sync.Mutex

func (conns *ActiveConnections) IsConn(userId uint) bool {
	mu.Lock()
	defer mu.Unlock()

	return len((*conns)[userId]) > 0
}

func (conns *ActiveConnections) Add(user *app.User, origin string, w http.ResponseWriter, r http.Request) *Conn {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("[%s] UserConn.Add TRACE user[%d] added from conn[%s]\n", h.GetReqId(&r), user.Id, origin)
	id := len(*conns)
	newConn := Conn{
		Id:     id,
		User:   user,
		Origin: origin,
		Writer: w,
		Reader: r,
		In:     make(chan event.LiveEvent, 64),
		//Out:    make(chan event.LiveUpdate, 64),
	}
	(*conns)[user.Id] = append((*conns)[user.Id], &newConn)
	return &newConn
}

func (conns *ActiveConnections) Get(userId uint) []*Conn {
	mu.Lock()
	defer mu.Unlock()
	return (*conns)[userId]
}

func (uc *ActiveConnections) Drop(c *Conn) error {
	mu.Lock()
	defer mu.Unlock()
	if c == nil || c.Origin == "" || c.User == nil {
		return fmt.Errorf("attempt to drop bad connection")
	}
	userConns := (*uc)[c.User.Id]
	for i, conn := range userConns {
		if conn.User == c.User && conn.Origin == c.Origin {
			(*uc)[c.User.Id] = append(userConns[:i], userConns[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("connection not found")
}

func (conn *Conn) SendUpdates(up event.LiveEvent, pollingUserId uint) {
	origin := conn.Origin
	if conn.User.Id != pollingUserId {
		log.Printf("[%s] Conn.SendUpdates WARN user[%v] is does not own conn[%v]\n", origin, pollingUserId, conn)
		return
	}
	err := conn.trySend(up)
	if err != nil {
		log.Printf("[%s] Conn.SendUpdates ERROR failed to send update to user[%d], err[%s]\n",
			origin, pollingUserId, err)
		up.Error = fmt.Errorf("ERROR SENDING TO: %d", pollingUserId)
		//conn.Out <- up
		return
	}
}

func (conn *Conn) trySend(up event.LiveEvent) error {
	w := conn.Writer
	if up.UserId <= 0 {
		return fmt.Errorf("Conn.trySend ERROR user is empty, user[%d], msg[%s]", up.UserId, up.Data)
	}
	if w == nil {
		return fmt.Errorf("Conn.trySend ERROR writer is nil")
	}
	switch up.Event {
	case event.UserChange,
		event.AvatarChange,
		event.ChatAdd,
		event.ChatInvite,
		event.ChatExpel,
		event.ChatLeave,
		event.ChatClose,
		event.ChatDrop,
		event.MessageAdd,
		event.MessageDrop:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to delete message to user[%d], %s", up.UserId, err)
		}
	default:
		return fmt.Errorf("Conn.trySend ERROR unknown update event[%s], update[%s]", up.Event, up.String())
	}
	return nil
}

func flushEvent(w *http.ResponseWriter, evnt event.EventType, up event.LiveEvent) error {
	if up.UserId <= 0 {
		panic("UserId should not be empty")
	}
	eventName := evnt.FormatEventName(up.ChatId, up.UserId, up.MsgId)
	eventId := utils.RandStringBytes(5)
	// must escape newlines in SSE
	data := utils.ReplaceWithSingleSpace(up.Data)
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

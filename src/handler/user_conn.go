package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

type Conn struct {
	Id     int
	User   *app.User
	Origin string
	Writer http.ResponseWriter
	Reader http.Request
	In     chan event.LiveUpdate // TODO test load
	//Out    chan event.LiveUpdate
}

type UserConn []Conn

var mu sync.Mutex

func (uc *UserConn) IsConn(userId uint) (bool, *Conn) {
	for _, conn := range *uc {
		if conn.User.Id == userId {
			return true, &conn
		}
	}
	return false, nil
}

func (uc *UserConn) Add(user *app.User, origin string, w http.ResponseWriter, r http.Request) *Conn {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("[%s] UserConn.Add TRACE user[%d] added from conn[%s]\n", h.GetReqId(&r), user.Id, origin)
	id := len(*uc)
	newConn := Conn{
		Id:     id,
		User:   user,
		Origin: origin,
		Writer: w,
		Reader: r,
		In:     make(chan event.LiveUpdate, 64),
		//Out:    make(chan event.LiveUpdate, 64),
	}
	*uc = append(*uc, newConn)
	return &newConn
}

func (uc UserConn) Get(userId uint) (*Conn, error) {
	mu.Lock()
	defer mu.Unlock()
	conns := uc.userConns(userId)
	if len(conns) <= 0 {
		return nil, fmt.Errorf("user[%d] not connected", userId)
	}

	log.Printf("UserConn.Get TRACE user[%d] has %d conns[%v]\n", userId, len(conns), conns)

	var conn *Conn
	for _, conn = range conns {
		if conn != nil && conn.User.Id == userId {
			break
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("user[%d] has no active conneciton", userId)
	}
	log.Printf("UserConn.Get TRACE user[%d] served on conn[%v]\n", userId, conn.Origin)
	return conn, nil
}

func (uc *UserConn) Drop(c *Conn) error {
	mu.Lock()
	defer mu.Unlock()
	if c == nil {
		return fmt.Errorf("attempt to drop NIL connection")
	}

	if uc == nil || len(*uc) <= 0 {
		return fmt.Errorf("no connections to drop")
	}

	for i, conn := range *uc {
		if conn.User == c.User && conn.Origin == c.Origin {
			*uc = append((*uc)[:i], (*uc)[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("connection not found")
}

func (uc *UserConn) userConns(userId uint) []*Conn {
	conns := make([]*Conn, 0)
	if uc == nil || len(*uc) <= 0 {
		return conns
	}
	for _, conn := range *uc {
		conn := conn
		if conn.User.Id == userId {
			conns = append(conns, &conn)
		}
	}
	return conns
}

func (conn *Conn) SendUpdates(up event.LiveUpdate, pollingUserId uint) {
	log.Printf("[%s] Conn.SendUpdates TRACE IN user[%d], input[%s]\n", conn.Origin, pollingUserId, up.String())
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
	log.Printf("[%s] Conn.SendUpdates TRACE OUT user[%d]\n", origin, pollingUserId)
}

func (conn *Conn) trySend(up event.LiveUpdate) error {
	w := conn.Writer
	if up.UserId <= 0 {
		return fmt.Errorf("Conn.trySend ERROR user is empty, user[%d], msg[%s]", up.UserId, up.Data)
	}
	if w == nil {
		return fmt.Errorf("Conn.trySend ERROR writer is nil")
	}
	switch up.Event {
	case event.UserChange:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to delete chat to user[%d], %s", up.UserId, err)
		}
	case event.AvatarChange:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to update avatar to user[%d], %s", up.UserId, err)
		}
	case event.ChatAdd, event.ChatInvite:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to send to user[%d], %s", up.UserId, err)
		}
	case event.ChatExpel:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to expel user[%d] from chat[%d], %s", up.UserId, up.ChatId, err)
		}
	case event.ChatLeave:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to leave user[%d] from chat[%d], %s", up.UserId, up.ChatId, err)
		}
	case event.ChatClose:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to close chat to user[%d], %s", up.UserId, err)
		}
	case event.ChatDrop:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to delete chat to user[%d], %s", up.UserId, err)
		}
	case event.MessageAdd:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to add message to user[%d], %s", up.UserId, err)
		}
	case event.MessageDrop:
		err := flushEvent(&w, up.Event, up)
		if err != nil {
			return fmt.Errorf("Conn.trySend ERROR failed to delete message to user[%d], %s", up.UserId, err)
		}
	default:
		return fmt.Errorf("Conn.trySend ERROR unknown update event[%v], update[%s]", up.Event, up.String())
	}
	return nil
}

func flushEvent(w *http.ResponseWriter, evnt event.UpdateType, up event.LiveUpdate) error {
	if up.UserId <= 0 {
		panic("UserId should not be empty")
	}
	eventName := evnt.FormatEventName(up.ChatId, up.UserId, up.MsgId)
	eventId := utils.RandStringBytes(5)
	// must escape newlines in SSE
	data := utils.TrimSpaces(up.Data)
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

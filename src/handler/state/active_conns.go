package state

import (
	"fmt"
	"log"
	"net/http"
	"prplchat/src/model/app"
	"prplchat/src/model/event"
	h "prplchat/src/utils/http"
	"sync"
)

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
	id := uint(len(*conns))
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

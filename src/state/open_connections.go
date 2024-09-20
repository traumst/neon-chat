package state

import (
	"fmt"
	"log"
	"neon-chat/src/app"
	"neon-chat/src/event"
	"net/http"
	"sync"
)

// TODO should be LRUCache
type OpenConnections map[uint][]*Conn

var mu sync.Mutex

func (conns *OpenConnections) IsConn(userId uint) bool {
	mu.Lock()
	defer mu.Unlock()

	return len((*conns)[userId]) > 0
}

func (conns *OpenConnections) Add(
	user *app.User,
	origin string,
	w http.ResponseWriter,
	r http.Request,
) *Conn {
	mu.Lock()
	defer mu.Unlock()
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
	log.Printf("UserConn.Add INFO added conn[%s] user[%d]\n", origin, user.Id)
	return &newConn
}

func (conns *OpenConnections) Get(userId uint) []*Conn {
	mu.Lock()
	defer mu.Unlock()
	return (*conns)[userId]
}

func (uc *OpenConnections) Drop(c *Conn) error {
	mu.Lock()
	defer mu.Unlock()
	if c == nil || c.Origin == "" || c.User == nil {
		return fmt.Errorf("attempt to drop bad connection")
	}
	userConns := (*uc)[c.User.Id]
	for i, conn := range userConns {
		if conn.User == c.User && conn.Origin == c.Origin {
			(*uc)[c.User.Id] = append(userConns[:i], userConns[i+1:]...)
			log.Printf("UserConn.Drop INFO dropped conn[%s] user[%d]\n", c.Origin, c.User.Id)
			return nil
		}
	}
	return fmt.Errorf("connection not found")
}

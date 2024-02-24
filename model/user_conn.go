package model

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/utils"
)

type Conn struct {
	ID     int
	User   string
	Origin string
	Writer http.ResponseWriter
	Reader http.Request
	In     chan UserUpdate
	Out    chan UserUpdate
}

type UserConn []Conn

var mu sync.Mutex

func (uc *UserConn) IsConn(user string) (bool, *Conn) {
	for _, conn := range *uc {
		if conn.User == user {
			return true, &conn
		}
	}
	return false, nil
}

func (uc *UserConn) Add(user string, origin string, w http.ResponseWriter, r http.Request) *Conn {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("∞---%s---> UserConn.Add TRACE user[%s] added from %s\n", utils.GetReqId(&r), user, origin)
	id := len(*uc)
	newConn := Conn{
		ID:     id,
		User:   user,
		Origin: origin,
		Writer: w,
		Reader: r,
		In:     make(chan UserUpdate, 64),
		Out:    make(chan UserUpdate, 64),
	}
	*uc = append(*uc, newConn)
	return &newConn
}

func (uc UserConn) Get(user string) (*Conn, error) {
	mu.Lock()
	defer mu.Unlock()
	conns := uc.userConns(user)
	if len(conns) == 0 {
		return nil, fmt.Errorf("user[%s] not connected", user)
	}

	log.Printf("∞--------> UserConn.Get TRACE user[%s] has %d conns[%v]\n", user, len(conns), conns)

	var conn *Conn
	for _, conn = range conns {
		if conn != nil && conn.User == user {
			break
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("user[%s] has no active conneciton", user)
	}
	log.Printf("∞--------> UserConn.Get TRACE user[%s] served on conn[%v]\n", user, conn)
	return conn, nil
}

func (uc *UserConn) Drop(c *Conn) error {
	mu.Lock()
	defer mu.Unlock()
	if c == nil {
		return fmt.Errorf("attempt to drop NIL connection")
	}

	if uc == nil || len(*uc) == 0 {
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

func (uc *UserConn) userConns(user string) []*Conn {
	conns := make([]*Conn, 0)
	if uc == nil || len(*uc) == 0 {
		return conns
	}
	for _, conn := range *uc {
		// TODO discard dead conn
		//userConn := conn // TODO this fixes distribution bug
		if conn.User == user {
			conns = append(conns, &conn)
			//conns = append(conns, &userConn)
		}
	}
	return conns
}

package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/utils"
)

type Conn struct {
	Id     int
	User   *app.User
	Origin string
	Writer http.ResponseWriter
	Reader http.Request
	In     chan event.LiveUpdate
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
	log.Printf("∞---%s---> UserConn.Add TRACE user[%d] added from conn[%s]\n", utils.GetReqId(&r), user.Id, origin)
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
	if len(conns) == 0 {
		return nil, fmt.Errorf("user[%d] not connected", userId)
	}

	log.Printf("∞--------> UserConn.Get TRACE user[%d] has %d conns[%v]\n", userId, len(conns), conns)

	var conn *Conn
	for _, conn = range conns {
		if conn != nil && conn.User.Id == userId {
			break
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("user[%d] has no active conneciton", userId)
	}
	log.Printf("∞--------> UserConn.Get TRACE user[%d] served on conn[%v]\n", userId, conn.Origin)
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

func (uc *UserConn) userConns(userId uint) []*Conn {
	conns := make([]*Conn, 0)
	if uc == nil || len(*uc) == 0 {
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

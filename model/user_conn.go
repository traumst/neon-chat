package model

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/utils"
)

type Conn struct {
	User    string
	Origin  string
	Writer  http.ResponseWriter
	Reader  http.Request
	Channel chan UserUpdate
}

func (c *Conn) Log() string {
	if c == nil {
		return "Conn: NIL"
	}
	return fmt.Sprintf("Conn:{User:\"%s\",Origin:\"%s\"}", c.User, c.Origin)
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
	newConn := Conn{
		User:    user,
		Origin:  origin,
		Writer:  w,
		Reader:  r,
		Channel: make(chan UserUpdate, 128),
	}
	*uc = append(*uc, newConn)
	return &newConn
}

func (uc UserConn) Get(reqId string, user string) (*Conn, error) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("∞---%s---> UserConn.Get TRACE IN user[%s]\n", reqId, user)
	conns := uc.userConns(reqId, user)
	if len(conns) == 0 {
		return nil, fmt.Errorf("user[%s] not connected", user)
	}
	for _, conn := range conns {
		if conn != nil {
			log.Printf("<---%s---∞ UserConn.Get TRACE OUT user[%s]\n", conn.Origin, user)
			return conn, nil
		}
	}

	log.Printf("<------∞ UserConn.Get ERROR OUT no conn to user[%s]\n", user)
	return nil, fmt.Errorf("user[%s] has no active conneciton", user)
}

func (uc *UserConn) Drop(reqId string, c *Conn) {
	mu.Lock()
	defer mu.Unlock()
	if c != nil {
		log.Printf("∞---%s---> UserConn.Drop TRACE user[%s]\n", reqId, c.User)
	} else {
		log.Printf("<---%s---∞ UserConn.Drop TRACE attempt to drop NIL connection\n", reqId)
		return
	}

	if uc == nil || len(*uc) == 0 {
		log.Printf("<---%s---∞ UserConn.Drop TRACE attempt to drop user[%s]\n", reqId, c.User)
		return
	}

	for i, conn := range *uc {
		if conn.User == c.User && conn.Origin == c.Origin {
			*uc = append((*uc)[:i], (*uc)[i+1:]...)
			return
		}
	}
}

func (uc *UserConn) userConns(reqId string, user string) []*Conn {
	conns := make([]*Conn, 0)
	if len(*uc) == 0 {
		log.Printf("<---%s---∞ UserConn.userConns TRACE user[%s] connections not found\n", reqId, user)
		return conns
	}
	for connID, conn := range *uc {
		if utils.GetReqId(&conn.Reader) == "" {
			log.Printf("∞---%s---∞ UserConn.userConns WARN user[%s] connection[%d] is NIL\n", reqId, user, connID)
			*uc = append((*uc)[:connID], (*uc)[connID+1:]...)
			continue
		}
		if conn.User == user {
			log.Printf("<---%s---∞ UserConn.userConns TRACE user[%s] connection[%d] found\n", reqId, user, connID)
			conns = append(conns, &conn)
		}
	}
	return conns
}

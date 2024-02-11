package model

import (
	"fmt"
	"log"
	"net/http"

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

func (uc *UserConn) userConns(user string) []*Conn {
	conns := make([]*Conn, 0)
	if uc == nil || len(*uc) == 0 {
		log.Printf("------ UserConn.userConns TRACE user[%s] connections not found\n", user)
		return conns
	}
	for connID, conn := range *uc {
		if conn.User == user {
			log.Printf("------ UserConn.userConns TRACE user[%s] connection[%d] found\n", user, connID)
			conns = append(conns, &conn)
		}
	}
	return conns
}

func (uc *UserConn) IsConn(user string) (bool, *Conn) {
	for _, conn := range *uc {
		if conn.User == user {
			return true, &conn
		}
	}
	return false, nil
}

func (uc *UserConn) Add(user string, origin string, w http.ResponseWriter, r http.Request) *Conn {
	log.Printf("--%s-- UserConn.Add TRACE user[%s] added from %s\n", utils.GetReqId(&r), user, origin)
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

func (uc *UserConn) Get(user string) (*Conn, error) {
	log.Printf("------ UserConn.Get TRACE user[%s]\n", user)
	conns := uc.userConns(user)
	if len(conns) == 0 {
		return nil, fmt.Errorf("user[%s] not connected", user)
	}
	for _, conn := range conns {
		if conn.Reader.Context().Done() != nil {
			return conn, nil
		} else {
			uc.Drop(conn)
		}
	}

	return nil, fmt.Errorf("user[%s] has no active conneciton", user)
}

func (uc *UserConn) Drop(c *Conn) {
	log.Printf("------ UserConn.Drop TRACE user[%s]\n", c.User)
	if uc == nil || len(*uc) == 0 {
		log.Printf("------ UserConn.Drop TRACE attempt to drop user[%s] origin[%s]\n", c.User, c.Origin)
		return
	}
	if c == nil {
		log.Printf("------ UserConn.Drop TRACE attempt to drop NIL connection\n")
		return
	}

	for i, conn := range *uc {
		if conn.User == c.User && conn.Origin == c.Origin {
			*uc = append((*uc)[:i], (*uc)[i+1:]...)
			return
		}
	}
}

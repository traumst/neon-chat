package model

import (
	"fmt"
	"log"
	"net/http"
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
	return fmt.Sprintf("Conn: User[%s] Origin[%s]", c.User, c.Origin)
}

type UserConn map[string]Conn

func (uc *UserConn) IsConn(user string) (bool, *Conn) {
	conn, ok := (*uc)[user]
	return ok, &conn
}

func (uc *UserConn) Add(user string, origin string, w http.ResponseWriter, r http.Request) *Conn {
	log.Printf("------ UserConn.Add TRACE user[%s] added from %s\n", user, origin)
	newConn := Conn{
		User:    user,
		Origin:  origin,
		Writer:  w,
		Reader:  r,
		Channel: make(chan UserUpdate, 128),
	}
	(*uc)[user] = newConn
	return &newConn
}

func (uc *UserConn) Get(user string) (*Conn, error) {
	log.Printf("------ UserConn.Get TRACE user[%s]\n", user)
	isConn, conn := uc.IsConn(user)
	if !isConn {
		return nil, fmt.Errorf("user not connected, %s", user)
	}
	return conn, nil
}

func (uc *UserConn) Drop(user string) {
	log.Printf("------ UserConn.Drop TRACE user[%s]\n", user)
	delete(*uc, user)
}

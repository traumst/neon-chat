package model

import (
	"fmt"
)

type UpdateType int

const (
	UnknownUpdate UpdateType = iota
	ChatUpdate    UpdateType = iota
	MessageUpdate UpdateType = iota
)

type UserUpdate struct {
	Type UpdateType
	User string
	Chat *Chat
	Msg  *Message
}

func (up *UserUpdate) Log() string {
	var chat string
	if up.Chat != nil {
		chat = up.Chat.Log()
	} else {
		chat = "nil"
	}

	var msg string
	if up.Msg != nil {
		msg = up.Msg.Log()
	} else {
		msg = "nil"
	}

	return fmt.Sprintf("UserUpdate{Type:%T,User:%s,Chat:%s,Msg:%s}",
		up.Type, up.User, chat, msg)
}

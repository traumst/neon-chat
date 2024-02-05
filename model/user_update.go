package model

import (
	"fmt"
)

type UpdateType int

const (
	ChatUpdate    UpdateType = iota
	MessageUpdate UpdateType = iota
)

type UserUpdate struct {
	Type UpdateType
	Chat *Chat
	Msg  *Message
	User string
}

func (up *UserUpdate) Log() string {
	return fmt.Sprintf("UserUpdate{Chat:%s,Msg:%s,User:%s}", up.Chat.Log(), up.Msg.Log(), up.User)
}

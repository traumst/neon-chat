package event

import "fmt"

type LiveUpdate struct {
	Event    UpdateType
	ChatId   int
	UserId   uint
	MsgId    int
	AuthorId uint
	Data     string
	Error    error
}

func (u *LiveUpdate) String() string {
	return fmt.Sprintf("LiveUpdate{event:%v,chat:%d,user:%d,msg:%d,author:%d,error:%v}",
		u.Event, u.ChatId, u.UserId, u.MsgId, u.AuthorId, u.Error)
}

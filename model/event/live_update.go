package event

import "fmt"

type LiveUpdate struct {
	Event  UpdateType
	ChatID int
	UserID string
	MsgID  int
	Author string
	Data   string
	Error  error
}

func (u *LiveUpdate) String() string {
	return fmt.Sprintf("LiveUpdate{event:%v,chat:%d,user:%s,msg:%d,author:%s,error:%v}",
		u.Event, u.ChatID, u.UserID, u.MsgID, u.Author, u.Error)
}

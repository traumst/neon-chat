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
	return fmt.Sprintf("LiveUpdate{event:%s,chat:%d,user:%s,msg:%d,author:%s,error:%v}",
		u.Event.String(), u.ChatID, u.UserID, u.MsgID, u.Author, u.Error)
}

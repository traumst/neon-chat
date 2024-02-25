package model

type UpdateType int
type UserUpdate struct {
	Type    UpdateType
	ChatID  int
	Author  string
	RawHtml string
}

const (
	UnknownUpdate UpdateType = iota
	ChatUpdate    UpdateType = iota
	ChatInvite    UpdateType = iota
	MessageUpdate UpdateType = iota
	ACK           UpdateType = iota
)

func (u UpdateType) String() string {
	switch u {
	case ChatUpdate:
		return string(ChatEventName)
	case ChatInvite:
		return string(ChatEventName)
	case MessageUpdate:
		return string(MessageEventName)
	default:
		return string(Unknown)
	}
}

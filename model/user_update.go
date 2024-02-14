package model

type UpdateType int

func (u UpdateType) String() string {
	switch u {
	case ChatUpdate:
		return "ChatUpdate"
	case ChatInvite:
		return "ChatInvite"
	case MessageUpdate:
		return "MessageUpdate"
	case PingUpdate:
		return "PingUpdate"
	default:
		return "UnknownUpdate"
	}
}

const (
	UnknownUpdate UpdateType = iota
	ChatUpdate    UpdateType = iota
	ChatInvite    UpdateType = iota
	MessageUpdate UpdateType = iota
	PingUpdate    UpdateType = iota
)

type UserUpdate struct {
	Type UpdateType
	User string
	Msg  string
}

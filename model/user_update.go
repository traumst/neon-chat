package model

type UpdateType int
type UserUpdate struct {
	Type   UpdateType
	ChatID int
	Author string
	Msg    string
}

const (
	UnknownUpdate UpdateType = iota
	ChatUpdate    UpdateType = iota
	ChatInvite    UpdateType = iota
	MessageUpdate UpdateType = iota
)

func (u UpdateType) String() string {
	switch u {
	case ChatUpdate:
		return "chat"
	case ChatInvite:
		return "invite"
	case MessageUpdate:
		return "message"
	default:
		return "unknown"
	}
}

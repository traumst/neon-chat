package model

type SSEvent string

const (
	Unknown          SSEvent = "unknown"
	PingEventName    SSEvent = "ping"
	ChatEventName    SSEvent = "chat"
	MessageEventName SSEvent = "msg"
	// TODO invite, msg in other chat
)

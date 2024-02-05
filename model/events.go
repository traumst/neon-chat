package model

type SSEvent string

const (
	PingEventName    SSEvent = "ping"
	ChatEventName    SSEvent = "chat"
	MessageEventName SSEvent = "message"
	// TODO invite, msg in other chat
)

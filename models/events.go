package models

type SSEvent string

const (
	ChatEventName    SSEvent = "chat"
	MessageEventName SSEvent = "message"
	// TODO invite, msg in other chat
)

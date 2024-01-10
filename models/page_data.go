package models

type PageData struct {
	Chats    []ChatCollapsed
	OpenChat *Chat
	Username string
}

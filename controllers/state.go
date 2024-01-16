package controllers

import (
	"net/http"

	"go.chat/models"
)

var chats = models.ChatList{}

func reqId(r *http.Request) string {
	return r.Header.Get("X-Request-Id")
}

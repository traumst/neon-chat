package controllers

import (
	"net/http"

	"go.chat/models"
	"go.chat/utils"
)

var chats = models.ChatList{}

func SetReqId(r *http.Request) string {
	reqId := utils.RandStringBytes(5)
	r.Header.Set("X-Request-Id", reqId)
	return reqId
}

func GetReqId(r *http.Request) string {
	return r.Header.Get("X-Request-Id")
}

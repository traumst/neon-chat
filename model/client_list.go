package model

import (
	"net/http"
)

type ClientList []*Client

type Client struct {
	User     string
	Request  *http.Request
	Response *http.ResponseWriter
}

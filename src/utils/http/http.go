package http

import (
	"net/http"

	"neon-chat/src/utils"
)

func SetReqId(r *http.Request, s *string) string {
	var reqId string
	if s == nil {
		reqId = utils.RandStringBytes(5)
	} else {
		reqId = *s
	}
	r.Header.Set("X-Request-Id", reqId)
	return reqId
}

func GetReqId(r *http.Request) string {
	return r.Header.Get("X-Request-Id")
}

func SetSseHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/event-stream")
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).Header().Set("Connection", "keep-alive")
}

func ParseUrlArgs(r *http.Request) map[string]string {
	args := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) <= 0 {
			args[k] = v[0]
		}
	}
	return args
}

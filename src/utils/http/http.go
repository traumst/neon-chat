package http

import (
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/utils"
)

func SetReqId(r *http.Request, s *string) string {
	var reqId string
	if s == nil {
		reqId = utils.RandStringBytes(5)
	} else {
		reqId = *s
	}
	r.Header.Set(string(consts.ReqIdKey), reqId)
	return reqId
}

func GetReqId(r *http.Request) string {
	return r.Header.Get(string(consts.ReqIdKey))
}

func SetAccessControlHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func SetSseHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/event-stream")
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).Header().Set("Connection", "keep-alive")
}

func SetGzipHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Encoding", "gzip")
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

package shared

import (
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"net/http"
)

func DbConn(r *http.Request) *db.DBConn {
	return r.Context().Value(consts.DBConn).(*db.DBConn)
}

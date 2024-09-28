package pub

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func ReadSession(
	state *state.State,
	dbConn *db.DBConn,
	w http.ResponseWriter,
	r *http.Request,
) (*app.User, error) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	cookie, err := h.GetSessionCookie(r)
	if err != nil {
		h.ClearSessionCookie(w, 0)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var appUser *app.User
	dbUser, err1 := db.GetUser(dbConn.Conn, cookie.UserId)
	if err1 != nil {
		h.ClearSessionCookie(w, 0)
		err = fmt.Errorf("failed to get user[%d] from cookie[%v], %s", cookie.UserId, cookie, err1.Error())
	} else {
		log.Printf("TRACE [%s] reading session user[%d][%s] from db\n", reqId, dbUser.Id, dbUser.Name)
		dbAvatar, _ := db.GetAvatar(dbConn.Conn, dbUser.Id)
		appUser = convert.UserDBToApp(dbUser, dbAvatar)
	}

	log.Printf("TRACE [%s] user has session:%t\n", reqId, err == nil)
	return appUser, err
}

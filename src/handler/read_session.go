package handler

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func ReadSession(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) (*a.User, error) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] ReadSession TRACE IN\n", reqId)
	cookie, err := h.GetSessionCookie(r)
	log.Printf("[%s] ReadSession TRACE session cookie[%v], err[%s]\n", reqId, cookie, err)
	if err != nil {
		h.ClearSessionCookie(w, 0)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var appUser *a.User
	dbUser, err1 := d.GetUser(db.Conn, cookie.UserId)
	if err1 != nil {
		h.ClearSessionCookie(w, 0)
		err = fmt.Errorf("failed to get user[%d] from cookie[%v], %s",
			cookie.UserId, cookie, err1.Error())
	} else {
		log.Printf("[%s] ReadSession TRACE session user[%d][%s], err[%s]\n",
			reqId, dbUser.Id, dbUser.Name, err1)
		dbAvatar, _ := d.GetAvatar(db.Conn, dbUser.Id)
		appUser = convert.UserDBToApp(dbUser, dbAvatar)
	}

	log.Printf("[%s] ReadSession TRACE OUT, success:%t\n", reqId, err == nil)
	return appUser, err
}

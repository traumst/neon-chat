package handler

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func ReadSession(state *state.State, dbConn *db.DBConn, w http.ResponseWriter, r *http.Request) (*app.User, error) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] ReadSession TRACE IN\n", reqId)
	cookie, err := h.GetSessionCookie(r)
	log.Printf("[%s] ReadSession TRACE session cookie[%v], err[%s]\n", reqId, cookie, err)
	if err != nil {
		h.ClearSessionCookie(w, 0)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var appUser *app.User
	dbUser, err1 := db.GetUser(dbConn.Conn, cookie.UserId)
	if err1 != nil {
		h.ClearSessionCookie(w, 0)
		err = fmt.Errorf("failed to get user[%d] from cookie[%v], %s",
			cookie.UserId, cookie, err1.Error())
	} else {
		log.Printf("[%s] ReadSession TRACE session user[%d][%s], err[%s]\n",
			reqId, dbUser.Id, dbUser.Name, err1)
		dbAvatar, _ := db.GetAvatar(dbConn.Conn, dbUser.Id)
		appUser = convert.UserDBToApp(dbUser, dbAvatar)
	}

	log.Printf("[%s] ReadSession TRACE OUT, success:%t\n", reqId, err == nil)
	return appUser, err
}

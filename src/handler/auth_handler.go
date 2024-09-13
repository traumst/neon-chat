package handler

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

func ReadSession(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) (*a.User, error) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
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

func Authenticate(
	db *d.DBConn,
	username string,
	pass string,
	authType a.AuthType,
) (*a.User, *a.Auth, error) {
	if db == nil || len(username) <= 0 || len(pass) <= 0 {
		log.Printf("Authenticate ERROR bad arguments username[%s] authType[%s]\n", username, authType)
		return nil, nil, fmt.Errorf("bad arguments")
	}
	dbUser, err := d.SearchUser(db.Conn, username)
	if err != nil || dbUser == nil || dbUser.Id <= 0 || len(dbUser.Salt) <= 0 {
		log.Printf("Authenticate TRACE user[%s] not found, result[%v], %s\n", username, dbUser, err)
		return nil, nil, nil
	}
	dbAvatar, _ := d.GetAvatar(db.Conn, dbUser.Id)
	appUser := convert.UserDBToApp(dbUser, dbAvatar)
	if appUser.Status != a.UserStatusActive {
		log.Printf("Authenticate WARN user[%d] status[%s] is inactive\n", dbUser.Id, dbUser.Status)
		return appUser, nil, nil
	}
	hash, err := utils.HashPassword(pass, appUser.Salt)
	if err != nil {
		log.Printf("Authenticate TRACE failed on hashing[%s] pass for user[%d], %s", hash, appUser.Id, err)
		return appUser, nil, fmt.Errorf("failed hashing pass for user[%d], %s", appUser.Id, err)
	}
	log.Printf("Authenticate TRACE user[%d] auth[%s] hash[%s]\n", appUser.Id, authType, hash)
	dbAuth, err := d.GetAuth(db.Conn, string(authType), hash)
	if err != nil {
		return appUser, nil, fmt.Errorf("no auth for user[%d] hash[%s], %s", appUser.Id, hash, err)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appUser, appAuth, nil
}

func Register(
	db *d.DBConn,
	newUser *a.User,
	pass string,
	authType a.AuthType,
) (*a.User, *a.Auth, error) {
	log.Printf("Register TRACE IN user\n")
	if db == nil || newUser == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", newUser, db)
	}
	if newUser.Name == "" || pass == "" || newUser.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] salt[%s]", newUser.Name, newUser.Salt)
	}
	var appUser *a.User
	var err error
	if newUser.Id != 0 {
		appUser = newUser
	} else {
		appUser, err = shared.CreateUser(db.Tx, newUser)
		if err != nil || appUser == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", appUser.Name)
		}
	}
	auth, err := shared.CreateAuth(db.Tx, appUser, pass, authType)
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, appUser, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", appUser.Id, auth)
	return appUser, auth, nil
}

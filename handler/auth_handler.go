package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/db"
	"go.chat/model"
	a "go.chat/model/app"
	"go.chat/utils"
)

func ReadSession(
	app *model.AppState,
	w http.ResponseWriter,
	r *http.Request,
) (*a.User, error) {
	log.Printf("--%s-> ReadSession TRACE IN\n", utils.GetReqId(r))
	cookie, err := utils.GetSessionCookie(r)
	log.Printf("--%s-> ReadSession TRACE session cookie[%v], err[%s]\n", utils.GetReqId(r), cookie, err)
	if err != nil {
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var user *a.User
	err = nil
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err = app.GetUser(cookie.UserId)
		log.Printf("--%s-> ReadSession TRACE session user[%v], err[%s]\n", utils.GetReqId(r), user, err)
		if user == nil {
			utils.ClearSessionCookie(w)
			err = fmt.Errorf("failed to get user from cookie[%v]", cookie)
		}
	}()
	wg.Wait()
	log.Printf("--%s-> ReadSession TRACE session user[%v], err[%s]\n", utils.GetReqId(r), user, err)
	return user, err
}

func Authenticate(
	db *db.DBConn,
	username string,
	pass string,
	authType a.AuthType,
) (*a.User, *a.UserAuth, error) {
	user, err := db.GetUser(username)
	if err != nil {
		return nil, nil, fmt.Errorf("username is already taken user[%s], %s", username, err)
	}
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed on hashing pass for user[%s], %s", username, err)
	}
	auth, err := db.GetAuth(user.Id, authType, *hash)
	if err != nil {
		return nil, nil, fmt.Errorf("no auth for user[%s], %s", username, err)
	}
	return user, auth, nil
}

func Register(
	db *db.DBConn,
	username string,
	pass string,
	authType a.AuthType,
) (*a.User, *a.UserAuth, error) {
	// TODO better salt
	salt := []byte(utils.RandStringBytes(16))
	user, err := db.AddUser(a.User{
		Id:   0,
		Name: username,
		Type: a.UserTypeFree,
		Salt: salt,
	})
	if err != nil || user == nil {
		return nil, nil, fmt.Errorf("failed to add user, %s", err)
	}
	if user.Id == 0 {
		return nil, nil, fmt.Errorf("user[%s] was not created", username)
	}
	hash, err := utils.HashPassword(pass, salt)
	if err != nil {
		return nil, nil, fmt.Errorf("error hashing pass, %s", err)
	}
	auth, err := db.AddAuth(a.UserAuth{
		Id:     0,
		UserId: user.Id,
		Type:   authType,
		Hash:   *hash,
	})
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("faild to add auth, %s", err)
	}
	if auth.Id == 0 {
		return nil, nil, fmt.Errorf("user[%s] auth was not created", username)
	}
	return user, auth, nil
}

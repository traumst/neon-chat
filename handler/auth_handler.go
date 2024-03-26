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
		return nil, nil, fmt.Errorf("user[%s] not found, %s", username, err)
	}
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		return user, nil, fmt.Errorf("failed on hashing pass for user[%s], %s", username, err)
	}
	log.Printf("-----> Authenticate TRACE user[%d] auth[%s] hash[%d]\n", user.Id, authType, *hash)
	auth, err := db.GetAuth(user.Id, authType, *hash)
	if err != nil {
		return user, nil, fmt.Errorf("no auth for user[%s] hash[%d], %s", username, hash, err)
	}
	return user, auth, nil
}

// Register creates a new user and its auth
// if user exists without auth - will only create the auth
func Register(
	db *db.DBConn,
	u *a.User,
	pass string,
	authType a.AuthType,
) (*a.User, *a.UserAuth, error) {
	log.Printf("-----> Register TRACE user[%v]\n", u)
	if db == nil || u == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", u, db)
	}
	// TODO think: forces to change salt when switching user.type
	seed := fmt.Sprintf("%s-%s", u.Name, u.Type)
	compl64 := utils.RandStringBytes(63 - len(seed))
	salt := fmt.Sprintf("%s;%s", seed, compl64)
	saltBytes := []byte(salt)
	var user *a.User
	var err error
	if u.Id == 0 {
		log.Printf("-----> Register TRACE creating user[%v]\n", u)
		user, err = db.AddUser(&a.User{
			Id:   0,
			Name: u.Name,
			Type: a.UserTypeFree,
			Salt: saltBytes,
		})
		if err != nil || user == nil {
			return nil, nil, fmt.Errorf("failed to add user, %s", err)
		}
		if user.Id == 0 {
			return nil, nil, fmt.Errorf("user[%s] was not created", user.Name)
		}
	}
	// TODO we may be in a partial state if we added user, but failed on auth
	hash, err := utils.HashPassword(pass, saltBytes)
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
		return nil, nil, fmt.Errorf("fail to add auth to user[%d][%s], %s", u.Id, u.Name, err)
	}
	if auth.Id == 0 {
		return nil, nil, fmt.Errorf("user[%d][%s] auth was not created", u.Id, u.Name)
	}
	return user, auth, nil
}

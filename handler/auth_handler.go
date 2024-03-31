package handler

import (
	"crypto/sha256"
	"encoding/hex"
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
		utils.ClearSessionCookie(w)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var user *a.User
	err = nil
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err = app.GetUser(cookie.UserId)
		if user == nil {
			utils.ClearSessionCookie(w)
			err = fmt.Errorf("failed to get user from cookie[%v]", cookie)
		} else {
			log.Printf("--%s-> ReadSession TRACE session user[%d][%s], err[%s]\n", utils.GetReqId(r),
				user.Id, user.Name, err)
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
	if db == nil || len(username) == 0 || len(pass) == 0 {
		return nil, nil, fmt.Errorf("missing mandatory args user[%s], authType[%s]", username, authType)
	}
	user, err := db.GetUser(username)
	if err != nil || user == nil || user.Id == 0 || len(user.Salt) == 0 {
		log.Printf("-----> Authenticate TRACE user[%v][%s] not found, %s\n", user, username, err)
		return nil, nil, fmt.Errorf("user[%s] not found, %s", username, err)
	}
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		log.Printf("-----> Authenticate TRACE failed on hashing[%s] pass for user[%s], %s", hash, user.Name, err)
		return user, nil, fmt.Errorf("failed on hashing pass for user[%s], %s", user.Name, err)
	}
	log.Printf("-----> Authenticate TRACE user[%d] auth[%s] hash[%s]\n", user.Id, authType, hash)
	auth, err := db.GetAuth(user.Id, authType, hash)
	if err != nil {
		return user, nil, fmt.Errorf("no auth for user[%s] hash[%s], %s", user.Name, hash, err)
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
	log.Printf("-----> Register TRACE IN user\n")
	if db == nil || u == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", u, db)
	}
	// TODO sterilize user input
	if u.Name == "" || pass == "" || u.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] pass[%s] salt[%s]", u.Name, pass, u.Salt)
	}
	var user *a.User
	if u.Id == 0 {
		user, err := createUser(db, u)
		if err != nil || user == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", u, err)
		}
	} else {
		user = u
	}
	auth, err := createAuth(db, user, pass, authType)
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, user, err)
	}
	return user, auth, nil
}

func GenerateSalt(userName string, userType a.UserType) string {
	// TODO think: type forces salt change when switching user.type
	seed := fmt.Sprintf("%s-%s", userType, userName)
	saltPlain := fmt.Sprintf("%s;%s", utils.RandStringBytes(7), seed)
	salt := sha256.Sum256([]byte(saltPlain))
	saltHex := hex.EncodeToString(salt[:])
	return saltHex
}

func createUser(db *db.DBConn, user *a.User) (*a.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("-----> createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("-----> createUser TRACE creating user[%s]\n", user.Name)
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err = db.AddUser(user)
	}()
	wg.Wait()
	if err != nil || user == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", user, err)
	}
	if user.Id == 0 {
		return nil, fmt.Errorf("user[%s] was not created", user.Name)
	}
	return user, err
}

func createAuth(db *db.DBConn, user *a.User, pass string, authType a.AuthType) (*a.UserAuth, error) {
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		return nil, fmt.Errorf("error hashing pass, %s", err)
	}
	log.Printf("-----> createAuth TRACE adding user[%d] auth[%s] hash[%s]\n", user.Id, authType, hash)
	var auth *a.UserAuth
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		auth, err = db.AddAuth(a.UserAuth{
			Id:     0,
			UserId: user.Id,
			Type:   authType,
			Hash:   hash,
		})
	}()
	wg.Wait()
	if err != nil || auth == nil {
		return nil, fmt.Errorf("fail to add auth to user[%d][%s], %s", user.Id, user.Name, err)
	}
	if auth.Id == 0 {
		return nil, fmt.Errorf("user[%d][%s] auth was not created", user.Id, user.Name)
	}
	return auth, err
}

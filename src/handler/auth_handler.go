package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"prplchat/src/convert"
	d "prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func ReadSession(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) (*a.User, error) {
	log.Printf("[%s] ReadSession TRACE IN\n", h.GetReqId(r))
	cookie, err := h.GetSessionCookie(r)
	log.Printf("[%s] ReadSession TRACE session cookie[%v], err[%s]\n", h.GetReqId(r), cookie, err)
	if err != nil {
		h.ClearSessionCookie(w, 0)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var appUser *a.User
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbUser, err1 := db.GetUser(cookie.UserId)
		if err1 != nil {
			h.ClearSessionCookie(w, 0)
			err = fmt.Errorf("failed to get user[%d] from cookie[%v], %s",
				cookie.UserId, cookie, err1.Error())
		} else {
			log.Printf("[%s] ReadSession TRACE session user[%d][%s], err[%s]\n",
				h.GetReqId(r), dbUser.Id, dbUser.Name, err1)
			appUser = convert.UserDBToApp(dbUser)
		}
	}()
	wg.Wait()
	log.Printf("[%s] ReadSession TRACE OUT, success:%t\n", h.GetReqId(r), err == nil)
	return appUser, err
}

func Authenticate(db *d.DBConn, username string, pass string, authType a.AuthType) (*a.User, *a.Auth, error) {
	if db == nil || len(username) <= 0 || len(pass) <= 0 {
		log.Printf("Authenticate ERROR bad arguments username[%s] authType[%s]\n", username, authType)
		return nil, nil, fmt.Errorf("bad arguments")
	}
	dbUser, err := db.SearchUser(username)
	if err != nil || dbUser == nil || dbUser.Id <= 0 || len(dbUser.Salt) <= 0 {
		log.Printf("Authenticate TRACE user[%s] not found, result[%v], %s\n", username, dbUser, err)
		return nil, nil, nil
	}
	appUser := convert.UserDBToApp(dbUser)
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
	dbAuth, err := db.GetAuth(string(authType), hash)
	if err != nil {
		return appUser, nil, fmt.Errorf("no auth for user[%d] hash[%s], %s", appUser.Id, hash, err)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appUser, appAuth, nil
}

func Register(db *d.DBConn, newUser *a.User, pass string, authType a.AuthType) (*a.User, *a.Auth, error) {
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
		appUser, err = createUser(db, newUser)
		if err != nil || appUser == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", appUser.Name)
		}
	}
	auth, err := createAuth(db, appUser, pass, authType)
	if err != nil || auth == nil {
		if recoverErr := deleteUser(db, appUser); recoverErr != nil {
			panic(fmt.Sprintf("failed to recovery-delete user[%d][%s], %s", appUser.Id, appUser.Name, err))
		}

		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, appUser, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", appUser.Id, auth)
	return appUser, auth, nil
}

func createUser(db *d.DBConn, user *a.User) (*a.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("createUser TRACE creating user[%s]\n", user.Name)
	dbUser := convert.UserAppToDB(user)
	created, err := db.AddUser(dbUser)
	if err != nil || created == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", created, err)
	}
	if created.Id <= 0 {
		return nil, fmt.Errorf("user[%s] was not created", created.Name)
	}
	appUser := convert.UserDBToApp(created)
	return appUser, err
}

func deleteUser(db *d.DBConn, user *a.User) error {
	if user.Id < 1 {
		log.Printf("deleteUser TRACE completing user[%s] signup\n", user.Name)
		return nil
	}
	log.Printf("deleteUser TRACE creating user[%s]\n", user.Name)
	err := db.DropUser(user.Id)
	if err != nil {
		return fmt.Errorf("failed to delete, %s", err)
	}
	return nil
}

func createAuth(db *d.DBConn, user *a.User, pass string, authType a.AuthType) (*a.Auth, error) {
	log.Printf("createAuth TRACE IN user[%d] auth[%s]\n", user.Id, authType)
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		return nil, fmt.Errorf("error hashing pass, %s", err)
	}
	log.Printf("createAuth TRACE adding user[%d] auth[%s] hash[%s]\n", user.Id, authType, hash)
	dbAuth := &d.Auth{
		Id:     0,
		UserId: user.Id,
		Type:   string(authType),
		Hash:   hash,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbAuth, err = db.AddAuth(*dbAuth)
	}()
	wg.Wait()
	if err != nil || dbAuth == nil {
		return nil, fmt.Errorf("fail to add auth to user[%d][%s], %s", user.Id, user.Name, err)
	}
	if dbAuth.Id <= 0 {
		return nil, fmt.Errorf("user[%d][%s] auth was not created", user.Id, user.Name)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appAuth, err
}

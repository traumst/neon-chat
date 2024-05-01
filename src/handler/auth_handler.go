package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/src/db"
	a "go.chat/src/model/app"
	"go.chat/src/utils"
)

func ReadSession(
	app *AppState,
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
	if err != nil {
		return nil, err
	}
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
		log.Printf("-----> Authenticate TRACE user[%s] not found, %s\n", username, err)
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
	if u.Name == "" || pass == "" || u.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] salt[%s]", u.Name, u.Salt)
	}
	var user *a.User
	var err error
	if u.Id != 0 {
		user = u
	} else {
		user, err = createUser(db, u)
		if err != nil || user == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", u, err)
		} else {
			log.Printf("-----> Register TRACE user[%s] created\n", user.Name)
		}
	}
	auth, err := createAuth(db, user, pass, authType)
	if err != nil || auth == nil {
		if recoverErr := deleteUser(db, user); recoverErr != nil {
			panic(fmt.Sprintf("failed to recovery-delete user[%d][%s], %s", user.Id, user.Name, err))
		}

		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, user, err)
	}
	log.Printf("-----> Register TRACE user[%v] auth[%v] created\n", user, auth)
	return user, auth, nil
}

func createUser(db *db.DBConn, user *a.User) (*a.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("-----> createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("-----> createUser TRACE creating user[%s]\n", user.Name)
	created, err := db.AddUser(user)
	if err != nil || created == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", created, err)
	}
	if created.Id == 0 {
		return nil, fmt.Errorf("user[%s] was not created", created.Name)
	}
	return created, err
}

func deleteUser(db *db.DBConn, user *a.User) error {
	if user.Id < 1 {
		log.Printf("-----> deleteUser TRACE completing user[%s] signup\n", user.Name)
		return nil
	}
	log.Printf("-----> deleteUser TRACE creating user[%s]\n", user.Name)
	err := db.DropUser(user.Id)
	if err != nil {
		return fmt.Errorf("failed to delete, %s", err)
	}
	return nil
}

func createAuth(db *db.DBConn, user *a.User, pass string, authType a.AuthType) (*a.UserAuth, error) {
	log.Printf("-----> createAuth TRACE IN user[%v] auth[%s]\n", user, authType)
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

package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	d "go.chat/src/db"
	a "go.chat/src/model/app"
	"go.chat/src/model/template/email"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

func ReadSession(
	app *AppState,
	db *d.DBConn,
	w http.ResponseWriter,
	r *http.Request,
) (*a.User, error) {
	log.Printf("[%s] ReadSession TRACE IN\n", h.GetReqId(r))
	cookie, err := h.GetSessionCookie(r)
	log.Printf("[%s] ReadSession TRACE session cookie[%v], err[%s]\n", h.GetReqId(r), cookie, err)
	if err != nil {
		h.ClearSessionCookie(w, 0)
		return nil, fmt.Errorf("failed to read session cookie, %s", err)
	}
	var user a.User
	err = nil
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbUser, err1 := db.GetUser(cookie.UserId)
		if err1 != nil {
			h.ClearSessionCookie(w, 0)
			err = fmt.Errorf("failed to get user[%d] from cookie[%v], %s", cookie.UserId, cookie, err1.Error())
		} else {
			log.Printf("[%s] ReadSession TRACE session user[%d][%s], err[%s]\n", h.GetReqId(r),
				dbUser.Id, dbUser.Name, err1)
			user = UserFromDB(*dbUser)
		}
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] ReadSession TRACE OUT\n", h.GetReqId(r))
	return &user, err
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
	user, err := db.SearchUser(username)
	if err != nil || user == nil || user.Id <= 0 || len(user.Salt) <= 0 {
		log.Printf("Authenticate TRACE user[%s] not found, %s\n", username, err)
		return nil, nil, nil
	}
	appUser := UserFromDB(*user)
	if appUser.Status != a.UserStatusActive {
		log.Printf("Authenticate WARN user[%d] status[%s] is inactive\n", user.Id, user.Status)
		return &appUser, nil, nil
	}
	hash, err := utils.HashPassword(pass, appUser.Salt)
	if err != nil {
		log.Printf("Authenticate TRACE failed on hashing[%s] pass for user[%d], %s", hash, appUser.Id, err)
		return &appUser, nil, fmt.Errorf("failed hashing pass for user[%d], %s", appUser.Id, err)
	}
	log.Printf("Authenticate TRACE user[%d] auth[%s] hash[%s]\n", appUser.Id, authType, hash)
	auth, err := db.GetAuth(string(authType), hash)
	if err != nil {
		return &appUser, nil, fmt.Errorf("no auth for user[%d] hash[%s], %s", appUser.Id, hash, err)
	}
	appAuth := AuthFromDB(*auth)
	return &appUser, &appAuth, nil
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
	var user *a.User
	var err error
	if newUser.Id != 0 {
		user = newUser
	} else {
		user, err = createUser(db, newUser)
		if err != nil || user == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", user.Name)
		}
	}
	auth, err := createAuth(db, user, pass, authType)
	if err != nil || auth == nil {
		if recoverErr := deleteUser(db, user); recoverErr != nil {
			panic(fmt.Sprintf("failed to recovery-delete user[%d][%s], %s", user.Id, user.Name, err))
		}

		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, user, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", user.Id, auth)
	return user, auth, nil
}

func IssueReservationToken(
	app *AppState,
	db *d.DBConn,
	user *a.User,
) (*email.VerifyEmailTemplate, error) {
	token := utils.RandStringBytes(16)
	expire := time.Now().Add(1 * time.Hour)
	reserve, err := reserve(db, user, token, expire)
	if err != nil {
		log.Printf("IssueReservationToken ERROR reserving[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("")
	}
	emailConfig := app.SmtpConfig()
	tmpl := email.VerifyEmailTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Token:       reserve.Token,
		//TokenExpire: reserve.Expire.Format(time.RFC3339),
		TokenExpire: reserve.Expire.Format(time.Stamp),
	}
	err = sendSignupCompletionEmail(tmpl, emailConfig.User, emailConfig.Pass)
	if err != nil {
		log.Printf("IssueReservationToken ERROR sending email from [%s] to [%s], %s\n",
			emailConfig.User, user.Email, err.Error())
		return nil, fmt.Errorf("failed to send email to[%s]", user.Email)
	}
	return &tmpl, nil
}

func reserve(
	db *d.DBConn,
	user *a.User,
	token string,
	expire time.Time,
) (*d.Reservation, error) {
	reserve := &d.Reservation{
		Id:     0,
		UserId: user.Id,
		Token:  token,
		Expire: expire,
	}
	reserve, err := db.AddReservation(*reserve)
	if err != nil {
		return nil, fmt.Errorf("reserve[%s] for user[%d], %s", token, user.Id, err)
	} else if reserve == nil {
		return nil, fmt.Errorf("reserve[%s] for user[%d] got reserve NIL", token, user.Id)
	} else if reserve.Id <= 0 {
		return nil, fmt.Errorf("reserve[%s] for user[%d] got reserve id 0", token, user.Id)
	}
	return reserve, nil
}

func createUser(db *d.DBConn, user *a.User) (*a.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("createUser TRACE creating user[%s]\n", user.Name)
	dbUser := UserToDB(*user)
	created, err := db.AddUser(&dbUser)
	if err != nil || created == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", created, err)
	}
	if created.Id <= 0 {
		return nil, fmt.Errorf("user[%s] was not created", created.Name)
	}
	appUser := UserFromDB(*created)
	return &appUser, err
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
	appAuth := AuthFromDB(*dbAuth)
	return &appAuth, err
}

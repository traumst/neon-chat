package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.chat/src/db"
	"go.chat/src/handler"
	a "go.chat/src/model/app"
	"go.chat/src/utils"
)

const (
	// TODO provide authType as form input
	authType = a.AuthTypeLocal
)

func Login(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login TRACE IN\n", utils.GetReqId(r))
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Bad Request"))
		return
	}
	u := r.FormValue("login-user")
	p := r.FormValue("login-pass")
	if u == "" || p == "" {
		log.Printf("--%s-> Login TRACE empty user[%s]", utils.GetReqId(r), u)
		RenderHome(app, w, r)
		return
	}
	log.Printf("--%s-> Login TRACE authentication check for user[%s] auth[%s]\n",
		utils.GetReqId(r), u, authType)
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		err = app.TrackUser(user)
		if err != nil {
			log.Printf("--%s-> Login WARN on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
		}
		utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	log.Printf("<-%s-- Login ERROR on authenticate[%s], %s\n", utils.GetReqId(r), u, err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Credentials not found for [%s:%s]", authType, u)))
	log.Printf("<-%s-- Login TRACE OUT\n", utils.GetReqId(r))
}

func SignUp(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> SignUp TRACE IN\n", utils.GetReqId(r))
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	u := r.FormValue("signup-user")
	p := r.FormValue("signup-pass")
	log.Printf("--%s-> SignUp TRACE authentication check for user[%s] auth[%s]\n", utils.GetReqId(r), u, authType)
	if u == "" || p == "" {
		RenderHome(app, w, r)
		return
	}
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		log.Printf("--%s-> SignUp TRACE signedIn instead of signUp user[%s], %s\n", utils.GetReqId(r), u, err)
		err := app.TrackUser(user)
		if err != nil {
			log.Printf("--%s-> SignUp ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
		}
		utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if user != nil {
		log.Printf("--%s-> SignUp ERROR name[%s] already taken by user[%d], %s\n",
			utils.GetReqId(r), u, user.Id, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to register, [%s] is already taken"))
		return
	}

	log.Printf("--%s-> SignUp TRACE register new user[%s], %s\n", utils.GetReqId(r), u, err)
	salt := utils.GenerateSalt(u, string(a.UserTypeFree))
	user = &a.User{
		Name: u,
		Type: a.UserTypeFree,
		Salt: salt,
	}
	user, auth, err = handler.Register(db, user, p, authType)
	if err != nil || user == nil || auth == nil {
		log.Printf("--%s-> SignUp ERROR on register user[%v], %s\n", utils.GetReqId(r), user, err)
		// TODO handler.Delete(db, user) - to remove partial data
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to register user [%s:%s]", a.UserTypeFree, u)))
		return
	}
	err = app.TrackUser(user)
	if err != nil {
		log.Printf("--%s-> SignUp ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
	}
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("--%s-> SignUp TRACE OUT\n", utils.GetReqId(r))
}

func Logout(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout TRACE \n", utils.GetReqId(r))
	utils.ClearSessionCookie(w)
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

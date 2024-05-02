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
	h "go.chat/src/utils/http"
)

const (
	authType = a.AuthTypeLocal
)

func Login(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login TRACE IN\n", h.GetReqId(r))
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}
	u := r.FormValue("login-user")
	u = utils.TrimSpaces(u)
	u = utils.TrimSpecial(u)
	p := r.FormValue("login-pass")
	p = utils.TrimSpaces(p)
	p = utils.TrimSpecial(p)
	if u == "" || p == "" || len(u) < 4 || len(p) < 4 {
		log.Printf("--%s-> Login TRACE empty user[%s]", h.GetReqId(r), u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad login credentials"))
		return
	}
	log.Printf("--%s-> Login TRACE authentication check for user[%s] auth[%s]\n",
		h.GetReqId(r), u, authType)
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		err = app.TrackUser(user)
		if err != nil {
			log.Printf("--%s-> Login WARN on track user[%d][%s], %s\n", h.GetReqId(r), user.Id, user.Name, err)
		}
		h.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	log.Printf("<-%s-- Login ERROR on authenticate[%s], %s\n", h.GetReqId(r), u, err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Credentials not found for [%s:%s]", authType, u)))
	log.Printf("<-%s-- Login TRACE OUT\n", h.GetReqId(r))
}

func SignUp(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> SignUp TRACE IN\n", h.GetReqId(r))
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	u := r.FormValue("signup-user")
	u = utils.TrimSpaces(u)
	u = utils.TrimSpecial(u)
	p := r.FormValue("signup-pass")
	p = utils.TrimSpaces(p)
	p = utils.TrimSpecial(p)
	log.Printf("--%s-> SignUp TRACE authentication check for user[%s] auth[%s]\n", h.GetReqId(r), u, authType)
	if u == "" || p == "" || len(u) < 4 || len(p) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad signup credentials"))
		return
	}
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		log.Printf("--%s-> SignUp TRACE signedIn instead of signUp user[%s], %s\n", h.GetReqId(r), u, err)
		err := app.TrackUser(user)
		if err != nil {
			log.Printf("--%s-> SignUp ERROR on track user[%d][%s], %s\n", h.GetReqId(r), user.Id, user.Name, err)
		}
		h.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if user != nil {
		log.Printf("--%s-> SignUp ERROR name[%s] already taken by user[%d], %s\n",
			h.GetReqId(r), u, user.Id, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to register, [%s] is already taken"))
		return
	}

	log.Printf("--%s-> SignUp TRACE register new user[%s], %s\n", h.GetReqId(r), u, err)
	salt := utils.GenerateSalt(u, string(a.UserTypeFree))
	user = &a.User{
		Name: u,
		Type: a.UserTypeFree,
		Salt: salt,
	}
	user, auth, err = handler.Register(db, user, p, authType)
	if err != nil || user == nil || auth == nil {
		log.Printf("--%s-> SignUp ERROR on register user[%v], %s\n", h.GetReqId(r), user, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to register user [%s:%s]", a.UserTypeFree, u)))
		return
	}
	err = app.TrackUser(user)
	if err != nil {
		log.Printf("--%s-> SignUp ERROR on track user[%d][%s], %s\n", h.GetReqId(r), user.Id, user.Name, err)
	}
	h.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("--%s-> SignUp TRACE OUT\n", h.GetReqId(r))
}

func Logout(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout TRACE \n", h.GetReqId(r))
	h.ClearSessionCookie(w)
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

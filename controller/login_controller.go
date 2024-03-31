package controller

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"go.chat/db"
	"go.chat/handler"
	"go.chat/model"
	a "go.chat/model/app"
	"go.chat/utils"
)

const (
	// TODO provide authType as form input
	authType = a.AuthTypeLocal
)

func Login(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login TRACE IN\n", utils.GetReqId(r))
	switch r.Method {
	case "GET":
		renderLogin(w, r)
	case "POST":
		signIn(app, conn, w, r)
	case "PUT":
		signUp(app, conn, w, r)
	default:
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout TRACE \n", utils.GetReqId(r))
	utils.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func renderLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> renderLogin\n", utils.GetReqId(r))
	loginTmpl, _ := template.ParseFiles("html/login_page.html")
	w.WriteHeader(http.StatusOK)
	loginTmpl.Execute(w, nil)
}

func signIn(app *model.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> signIn TRACE IN\n", utils.GetReqId(r))
	u := r.FormValue("user")
	p := r.FormValue("pass")
	if u == "" || p == "" {
		renderLogin(w, r)
		return
	}
	log.Printf("--%s-> signIn TRACE authentication check for user[%s] auth[%s]\n",
		utils.GetReqId(r), u, authType)
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user == nil {
		log.Printf("--%s-> signIn INFO user[%s] not found\n", utils.GetReqId(r), u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found"))
		return
	}
	if auth == nil {
		log.Printf("--%s-> signIn INFO user[%s] has no auth\n", utils.GetReqId(r), u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Found"))
		return
	}
	if err != nil {
		log.Printf("--%s-> signIn ERROR on authenticate[%s], %s\n", utils.GetReqId(r), u, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Found"))
		return
	}

	err = app.TrackUser(user)
	if err != nil {
		log.Printf("--%s-> signIn ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
	}
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	renderLogin(w, r)
}

func signUp(app *model.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> signUp TRACE IN\n", utils.GetReqId(r))
	u := r.FormValue("user")
	p := r.FormValue("pass")
	if u == "" || p == "" {
		renderLogin(w, r)
		return
	}
	log.Printf("--%s-> signUp TRACE authentication check for user[%s] auth[%s]\n",
		utils.GetReqId(r), u, authType)
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		err := app.TrackUser(user)
		if err != nil {
			log.Printf("--%s-> signUp ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
		}
		utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if user != nil {
		log.Printf("--%s-> signUp TRACE completing user[%s], %s\n", utils.GetReqId(r), u, err)
	} else {
		log.Printf("--%s-> signUp TRACE register new user[%s], %s\n", utils.GetReqId(r), u, err)
		user = &a.User{
			Name: u,
			Type: a.UserTypeFree,
			Salt: handler.GenerateSalt(u, a.UserTypeFree),
		}
	}
	if user.Salt == "" {
		panic("user salt is empty")
	}
	user, auth, err = handler.Register(db, user, p, authType)
	if err != nil {
		log.Printf("--%s-> signUp ERROR on register user[%v], %s\n", utils.GetReqId(r), user, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Operation failed"))
		return
	}
	err = app.TrackUser(user)
	if err != nil {
		log.Printf("--%s-> signUp ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
	}
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("--%s-> signUp TRACE OUT\n", utils.GetReqId(r))
}

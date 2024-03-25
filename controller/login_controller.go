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

func Login(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login TRACE IN\n", utils.GetReqId(r))
	cookie, _ := utils.GetSessionCookie(r)
	if cookie != nil {
		user, _ := app.GetUser(cookie.UserId)
		if user != nil {
			log.Printf("--%s-> Login WARN user[%d] is already logged in, redirected\n", utils.GetReqId(r), cookie.UserId)
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		utils.ClearSessionCookie(w)
	}
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
	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
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
	if u == "" {
		log.Printf("--%s-> signIn ERROR user\n", utils.GetReqId(r))
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}
	p := r.FormValue("pass")
	if p == "" {
		log.Printf("--%s-> signIn ERROR pass\n", utils.GetReqId(r))
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	// TODO consider other auth types
	user, auth, err := handler.Authenticate(db, u, p, a.AuthTypeLocal)
	if err != nil {
		log.Printf("--%s-> signIn ERROR on authenticate[%s], %s\n", utils.GetReqId(r), u, err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_ = app.TrackUser(user)
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("<-%s-- signIn TRACE OUT\n", utils.GetReqId(r))
}

func signUp(app *model.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> signUp TRACE IN\n", utils.GetReqId(r))
	u := r.FormValue("user")
	if u == "" {
		log.Printf("--%s-> signUp ERROR user\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	p := r.FormValue("pass")
	if p == "" {
		log.Printf("--%s-> signUp ERROR pass\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	// TODO consider other auth types
	user, auth, _ := handler.Authenticate(db, u, p, a.AuthTypeLocal)
	if user != nil && auth != nil {
		log.Printf("--%s-> signUp WARN user[%s] already has auth[%s]\n", utils.GetReqId(r), u, a.AuthTypeLocal)
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}
	user, auth, err := handler.Register(db, u, p, a.AuthTypeLocal)
	if err != nil {
		log.Printf("--%s-> signUp ERROR on register user[%s], %s\n", utils.GetReqId(r), u, err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	_ = app.TrackUser(user)
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("--%s-> signUp TRACE OUT\n", utils.GetReqId(r))
}

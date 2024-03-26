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
	user, auth, err := handler.Authenticate(db, u, p, authType)
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
	p := r.FormValue("pass")
	if u == "" || p == "" {
		renderLogin(w, r)
		return
	}
	log.Printf("--%s-> signUp TRACE authentication check...\n", utils.GetReqId(r))
	user, auth, _ := handler.Authenticate(db, u, p, authType)
	if user != nil && auth != nil {
		log.Printf("--%s-> signUp INFO authenticated user[%s] on auth[%s]\n",
			utils.GetReqId(r), user.Name, auth.Type)
		utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	} else if user != nil {
		log.Printf("--%s-> signUp INFO user[%v] is still partial\n", utils.GetReqId(r), user)
		user = &a.User{Name: u, Type: a.UserTypeFree}
	} else {
		log.Printf("--%s-> signUp INFO user[%s] has no auth\n", utils.GetReqId(r), u)
		user = &a.User{Name: u, Type: a.UserTypeFree}
	}
	log.Printf("--%s-> signUp TRACE registrastion for user[%v]\n", utils.GetReqId(r), user)
	user, auth, err := handler.Register(db, user, u, authType)
	if err != nil {
		log.Printf("--%s-> signUp ERROR on register user[%v], %s\n", utils.GetReqId(r), user, err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	_ = app.TrackUser(user)
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("--%s-> signUp TRACE OUT\n", utils.GetReqId(r))
}

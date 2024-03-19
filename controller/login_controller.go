package controller

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"go.chat/db"
	"go.chat/model"
	"go.chat/utils"
)

// TODO WIP conn *db.DBConn
func Login(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login\n", utils.GetReqId(r))
	switch r.Method {
	case "GET":
		renderLogin(w, r)
	case "POST":
		signIn(w, r)
	default:
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func Logout(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout\n", utils.GetReqId(r))
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func renderLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> RenderLogin\n", utils.GetReqId(r))
	loginTmpl, _ := template.ParseFiles("html/login_page.html")
	loginTmpl.Execute(w, nil)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> SignIn\n", utils.GetReqId(r))
	user := r.FormValue("username")
	if user == "" {
		log.Printf("--%s-> SignIn ERROR username\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   user,
		Expires: time.Now().Add(8 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

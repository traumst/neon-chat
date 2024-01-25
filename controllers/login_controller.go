package controllers

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login\n", GetReqId(r))
	switch r.Method {
	case "GET":
		RenderLogin(w, r)
	case "POST":
		SignIn(w, r)
	default:
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> SignIn\n", GetReqId(r))
	username := r.FormValue("username")
	if username == "" {
		log.Printf("--%s-> SignIn ERROR username\n", GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   username,
		Expires: time.Now().Add(8 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> RenderLogin\n", GetReqId(r))
	loginTmpl, _ := template.ParseFiles("views/login.html")
	loginTmpl.Execute(w, nil)
}

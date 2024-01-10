package controllers

import (
	"html/template"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		if username != "" {
			// Set a cookie with the username as value
			http.SetCookie(w, &http.Cookie{
				Name:    "username",
				Value:   username,
				Expires: time.Now().Add(8 * time.Hour),
			})

			// Redirect to the chat page
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	loginTmpl, _ := template.ParseFiles("views/login.html")
	loginTmpl.Execute(w, nil)
}

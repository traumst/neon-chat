package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.chat/src"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

	// If not a POST request or no username, show the login page
	loginTmpl, _ := template.ParseFiles("pages/login.html")
	loginTmpl.Execute(w, nil)
}

var messageStore = src.MessageStore{}

type PageData struct {
	Messages []src.Message
	Username string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := PageData{
		Messages: messageStore.Get(),
		Username: usernameCookie.Value,
	}

	tmpl := template.Must(template.ParseFiles("pages/chat.html"))
	tmpl.Execute(w, data)
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		author := usernameCookie.Value
		text := template.HTMLEscapeString(r.FormValue("text"))
		if text != "" && author != "" {
			messageStore.Add(src.Message{Author: author, Text: text})
		}
	}
	// Redirect back to the main page
	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err == nil {
			messageStore.Delete(id)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	log.Println("Router setup...")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

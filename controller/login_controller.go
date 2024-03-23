package controller

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"go.chat/db"
	"go.chat/model"
	"go.chat/model/app"
	"go.chat/utils"
)

func Login(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login\n", utils.GetReqId(r))
	switch r.Method {
	case "GET":
		renderLogin(w, r)
	case "POST":
		signIn(w, r, conn)
	case "PUT":
		signUp(w, r, conn)
	default:
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func Logout(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout\n", utils.GetReqId(r))
	utils.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func renderLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> renderLogin\n", utils.GetReqId(r))
	cookie, err := utils.GetSessionCookie(r)
	if err == nil && cookie != nil {
		log.Printf("--%s-> renderLogin TRACE user already has cookie, redirecting to HOME\n", utils.GetReqId(r))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	loginTmpl, _ := template.ParseFiles("html/login_page.html")
	loginTmpl.Execute(w, nil)
}

func signIn(w http.ResponseWriter, r *http.Request, conn *db.DBConn) {
	log.Printf("--%s-> signIn\n", utils.GetReqId(r))
	cookie, err := utils.GetSessionCookie(r)
	if err == nil && cookie != nil {
		log.Printf("--%s-> signIn TRACE user already has cookie, redirecting to HOME\n", utils.GetReqId(r))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	username := r.FormValue("user")
	if username == "" {
		log.Printf("--%s-> signIn ERROR user\n", utils.GetReqId(r))
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}
	pass := r.FormValue("pass")
	if pass == "" {
		log.Printf("--%s-> signIn ERROR pass\n", utils.GetReqId(r))
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	user, err := conn.GetUser(username)
	if err != nil {
		log.Printf("--%s-> signIn ERROR on user[%s], %s\n", utils.GetReqId(r), username, err.Error())
		// TODO do not redirect
		http.Error(w, "Invalid username", http.StatusNotFound)
		return
	}
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		log.Printf("--%s-> signIn ERROR on hashing[%s], %s\n", utils.GetReqId(r), username, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	auth, err := conn.GetAuth(user.Id, app.AuthTypeLocal, *hash)
	if err != nil {
		log.Printf("--%s-> signIn ERROR name is already taken [%s]\n", utils.GetReqId(r), username)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
}

func signUp(w http.ResponseWriter, r *http.Request, conn *db.DBConn) {
	log.Printf("--%s-> signUp\n", utils.GetReqId(r))
	cookie, err := utils.GetSessionCookie(r)
	if err == nil && cookie != nil {
		log.Printf("--%s-> signUp TRACE user already has cookie, redirecting to HOME\n", utils.GetReqId(r))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	username := r.FormValue("user")
	if username == "" {
		log.Printf("--%s-> signUp ERROR user\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> signUp 2\n", utils.GetReqId(r))
	pass := r.FormValue("pass")
	if pass == "" {
		log.Printf("--%s-> signUp ERROR pass\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> signUp 3\n", utils.GetReqId(r))
	user, _ := conn.GetUser(username)
	if user != nil {
		// check if user completed signup
		log.Printf("--%s-> signUp 4\n", utils.GetReqId(r))
		hash, err := utils.HashPassword(pass, user.Salt)
		if err != nil {
			log.Printf("--%s-> signUp ERROR on hashing[%s], %s\n", utils.GetReqId(r), username, err.Error())
			http.Redirect(w, r, "/login", http.StatusInternalServerError)
			return
		}
		log.Printf("--%s-> signUp 5\n", utils.GetReqId(r))
		auth, _ := conn.GetAuth(user.Id, app.AuthTypeLocal, *hash)
		if auth != nil {
			log.Printf("--%s-> signUp ERROR name is already taken [%s]\n", utils.GetReqId(r), username)
			http.Redirect(w, r, "/login", http.StatusBadRequest)
			return
		}
	}
	log.Printf("--%s-> signUp 6\n", utils.GetReqId(r))
	// TODO better salt
	salt := []byte(utils.RandStringBytes(16))
	hash, err := utils.HashPassword(pass, salt)
	if err != nil {
		log.Printf("--%s-> signUp ERROR on hash[%s], %s\n", utils.GetReqId(r), username, err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	log.Printf("--%s-> signUp 7\n", utils.GetReqId(r))
	user, err = conn.AddUser(app.User{
		Id:   0,
		Name: username,
		Type: app.Free,
		Salt: salt,
	})
	if err != nil || user == nil {
		log.Printf("--%s-> signUp ERROR on user[%s], %s\n", utils.GetReqId(r), username, err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	if user.Id == 0 {
		log.Printf("--%s-> signUp ERROR user.Id is 0\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> signUp 8\n", utils.GetReqId(r))
	auth, err := conn.AddAuth(app.UserAuth{
		Id:     0,
		UserId: user.Id,
		Type:   app.AuthTypeLocal,
		Hash:   *hash,
	})
	if err != nil || auth == nil {
		log.Printf("--%s-> signUp ERROR on auth[%s], %s\n", utils.GetReqId(r), username, err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	if auth.Id == 0 {
		log.Printf("--%s-> signUp ERROR auth.Id is 0\n", utils.GetReqId(r))
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	log.Printf("--%s-> signUp 8\n", utils.GetReqId(r))
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	http.Redirect(w, r, "/", http.StatusFound)
}

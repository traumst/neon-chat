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

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> RenderLogin\n", utils.GetReqId(r))
	loginTmpl, _ := template.ParseFiles("html/login_page.html")
	w.WriteHeader(http.StatusOK)
	loginTmpl.Execute(w, nil)
}

func Login(app *model.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Login TRACE IN\n", utils.GetReqId(r))
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Bad Request"))
		return
	}
	u := r.FormValue("user")
	p := r.FormValue("pass")
	if u == "" || p == "" {
		RenderLogin(w, r)
		return
	}
	log.Printf("--%s-> Login TRACE authentication check for user[%s] auth[%s]\n",
		utils.GetReqId(r), u, authType)
	user, auth, err := handler.Authenticate(db, u, p, authType)
	if user == nil {
		log.Printf("--%s-> Login INFO user[%s] not found\n", utils.GetReqId(r), u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found"))
		return
	}
	if auth == nil {
		log.Printf("--%s-> Login INFO user[%s] has no auth\n", utils.GetReqId(r), u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Found"))
		return
	}
	if err != nil {
		log.Printf("--%s-> Login ERROR on authenticate[%s], %s\n", utils.GetReqId(r), u, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Found"))
		return
	}

	err = app.TrackUser(user)
	if err != nil {
		log.Printf("--%s-> Login ERROR on track user[%d][%s], %s\n", utils.GetReqId(r), user.Id, user.Name, err)
	}
	utils.SetSessionCookie(w, user, auth, time.Now().Add(8*time.Hour))
	RenderLogin(w, r)
}

func SignUp(app *model.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> SignUp TRACE IN\n", utils.GetReqId(r))
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Bad Request"))
		return
	}
	u := r.FormValue("user")
	p := r.FormValue("pass")
	if u == "" || p == "" {
		RenderLogin(w, r)
		return
	}
	log.Printf("--%s-> SignUp TRACE authentication check for user[%s] auth[%s]\n",
		utils.GetReqId(r), u, authType)
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
		log.Printf("--%s-> SignUp ERROR user[%s] already taken by user[%d], %s\n",
			utils.GetReqId(r), user.Name, user.Id, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Operation failed"))
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
		// TODO handler.Delete(db, user)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Operation failed"))
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

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Logout TRACE \n", utils.GetReqId(r))
	utils.ClearSessionCookie(w)
	RenderLogin(w, r)
}

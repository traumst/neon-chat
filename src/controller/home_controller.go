package controller

import (
	"log"
	"net/http"

	"prplchat/src/controller/shared"
	"prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func RenderLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("[%s] RenderLogin TRACE IN", h.GetReqId(r))
	login := template.AuthTemplate{}
	home := template.HomeTemplate{
		Chats:         nil,
		OpenChat:      nil,
		User:          template.UserTemplate{UserName: "anon"},
		IsAuthorized:  false,
		LoginTemplate: login,
		Avatar:        nil,
	}
	log.Printf("[%s] RenderLogin TRACE templating", h.GetReqId(r))
	html, err := home.HTML()
	if err != nil {
		log.Printf("[%s] RenderLogin ERROR login %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home login"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func RenderHome(
	state *state.State,
	db *db.DBConn,
	w http.ResponseWriter,
	r *http.Request,
	user *a.User,
) {
	if state == nil {
		panic("app is nil")
	} else if db == nil {
		panic("db is nil")
	} else if user == nil {
		panic("user is nil")
	}
	log.Printf("[%s] RenderHome TRACE IN", h.GetReqId(r))
	html, err := shared.TemplateHome(state, db, r, user)
	if err != nil {
		log.Printf("[%s] RenderHome ERROR failed to template home, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	log.Printf("[%s] RenderHome TRACE, user[%d] gets content\n", h.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

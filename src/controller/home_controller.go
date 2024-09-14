package controller

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	t "neon-chat/src/model/template"
)

func NavigateHome(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] NavigateHome\n", reqId)
	user := r.Context().Value(consts.ActiveUser)
	if user != nil {
		RenderHome(w, r)
	} else {
		RenderLogin(w, r)
	}
}

func RenderHome(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] RenderHome\n", reqId)
	ctx := r.Context()
	html, err := handler.TemplateHome(
		ctx.Value(consts.AppState).(*state.State),
		ctx.Value(consts.DBConn).(*d.DBConn),
		ctx.Value(consts.ActiveUser).(*a.User),
	)
	if err != nil {
		log.Printf("[%s] RenderHome ERROR failed to template home, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	login := t.AuthTemplate{}
	home := t.HomeTemplate{
		Chats:         nil,
		OpenChat:      nil,
		User:          t.UserTemplate{UserName: "anon"},
		IsAuthorized:  false,
		LoginTemplate: login,
		Avatar:        nil,
	}
	html, err := home.HTML()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home login"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

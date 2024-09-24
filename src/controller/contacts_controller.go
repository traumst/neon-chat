package controller

import (
	"log"
	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/controller/shared"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/parse"
	"neon-chat/src/template"
	"net/http"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] GetUserInfo\n", reqId)
	if r.Method != "GET" {
		log.Printf("TRACE [%s] GetUserInfo auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	args, err := parse.ParseQueryString(r)
	if err != nil {
		log.Printf("ERROR [%s] GetUserInfo parsing arguments, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	if args.UserId < 1 {
		log.Printf("TRACE [%s] GetUserInfo invalid user id\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid user id"))
		return
	}
	viewer := r.Context().Value(consts.ActiveUser).(*app.User)
	dbConn := shared.DbConn(r)
	dbUser, err := db.GetUser(dbConn.Conn, args.UserId)
	if err != nil {
		log.Printf("[%s] GetUserInfo ERROR retrieving user[%d] data %s\n", reqId, args.UserId, err)
	}
	if dbUser == nil {
		log.Printf("[%s] GetUserInfo TRACE user[%d] not found\n", reqId, args.UserId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	dbAvatar, err := db.GetAvatar(dbConn.Conn, dbUser.Id)
	if err != nil {
		log.Printf("[%s] GetUserInfo ERROR retrieving user[%d] avatar: %s\n", reqId, dbUser.Id, err)
	}
	appUser := convert.UserDBToApp(dbUser, dbAvatar)
	avatar := appUser.Avatar.Template(viewer)
	tmpl := template.UserInfoTemplate{
		ViewerId:     viewer.Id,
		UserId:       appUser.Id,
		UserName:     appUser.Name,
		UserEmail:    appUser.Email,
		UserAvatar:   avatar,
		RegisterDate: "",
	}
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("[%s] GetUserInfo ERROR cannot template response: %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

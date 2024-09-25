package controller

import (
	"fmt"
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
	tmpl, err := GetUserInfoCard(w, r, viewer, args.UserId)
	if err != nil {
		log.Printf("[%s] GetUserInfo TRACE user[%d] not found\n", reqId, args.UserId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
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

func GetUserInfoCard(
	w http.ResponseWriter,
	r *http.Request,
	viewer *app.User,
	otherUserId uint,
) (template.UserInfoTemplate, error) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	dbConn := shared.DbConn(r)
	dbUser, err := db.GetUser(dbConn.Conn, otherUserId)
	if err != nil {
		log.Printf("[%s] GetUserInfo ERROR retrieving user[%d] data %s\n", reqId, otherUserId, err)
	}
	if dbUser == nil {
		log.Printf("[%s] GetUserInfo TRACE user[%d] not found\n", reqId, otherUserId)
		return template.UserInfoTemplate{}, fmt.Errorf("user not found")
	}
	dbAvatar, err := db.GetAvatar(dbConn.Conn, dbUser.Id)
	if err != nil {
		log.Printf("[%s] GetUserInfo ERROR retrieving user[%d] avatar: %s\n", reqId, dbUser.Id, err)
	}
	appUser := convert.UserDBToApp(dbUser, dbAvatar)
	avatar := appUser.Avatar.Template(viewer)
	sharedChats := GetSharedChats(dbConn, viewer, dbUser)
	tmpl := template.UserInfoTemplate{
		ViewerId:    viewer.Id,
		UserId:      appUser.Id,
		UserName:    appUser.Name,
		UserEmail:   appUser.Email,
		UserAvatar:  avatar,
		SharedChats: sharedChats,
	}
	return tmpl, nil
}

func GetSharedChats(dbConn *db.DBConn, viewer *app.User, otherUser *db.User) []template.ChatTemplate {
	dbChats, _ := db.GetSharedChats(dbConn.Conn, []uint{viewer.Id, otherUser.Id})
	if dbChats == nil {
		dbChats = []db.Chat{}
	}
	sharedChats := make([]template.ChatTemplate, 0, len(dbChats))
	for _, dbChat := range dbChats {
		shared := convert.ChatDBToApp(&dbChat, &db.User{
			Id:     viewer.Id,
			Name:   viewer.Name,
			Email:  viewer.Email,
			Type:   string(viewer.Type),
			Status: string(viewer.Status),
			Salt:   "",
		})
		sharedChats = append(sharedChats, template.ChatTemplate{
			ChatId:   shared.Id,
			ChatName: shared.Name,
		})
	}
	return sharedChats
}

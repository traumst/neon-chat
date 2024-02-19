package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.chat/model"
	"go.chat/utils"
)

type ChatController struct {
}

func (c *ChatController) OpenChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> OpenChat\n", utils.GetReqId(r))
	if r.Method != "GET" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", utils.GetReqId(r), path[2])
	chatID, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if chatID < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", utils.GetReqId(r), chatID)
	openChat, err := app.State.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", utils.GetReqId(r))
	html, err := openChat.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- OpenChat TRACE returning template\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func (c *ChatController) AddChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddChat\n", utils.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- AddChat TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("--%s-> AddChat TRACE check login\n", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	chatName := r.FormValue("chatName")
	log.Printf("--%s-> AddChat TRACE adding user[%s] chat[%s]\n", utils.GetReqId(r), user, chatName)
	chatID := app.State.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE user[%s] opening chat[%s][%d]\n", utils.GetReqId(r), user, chatName, chatID)
	openChat, err := app.State.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", utils.GetReqId(r), err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}

	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n", utils.GetReqId(r), chatName, chatID)
	template := openChat.ToTemplate(user)
	sendChatContent(utils.GetReqId(r), w, template)
	err = sendChatHeader(utils.GetReqId(r), template)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR cannot distribute chat header, %s\n", utils.GetReqId(r), err)
	}
}

func sendChatContent(reqId string, w http.ResponseWriter, template *model.ChatTemplate) {
	html, err := template.GetHTML()
	if err != nil {
		log.Printf("<--%s-- sendChatContent ERROR cannot template chat [%+v], %s", reqId, template, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- sendChatContent TRACE writing response\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func sendChatHeader(reqId string, template *model.ChatTemplate) error {
	shortHtml, err := template.GetShortHTML()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errors := make([]string, 0)
	log.Printf("--%s-> sendChatHeader TRACE distributing chat[%s] header to users [%+v]\n",
		reqId, template.Name, template.Users)
	for _, user := range template.Users {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			conn, err := app.State.GetConn(user)
			if err != nil {
				errors = append(errors, "user:"+user+",err:"+err.Error())
				return
			}
			log.Printf("--%s-> sendChatHeader TRACE distributing chat[%s] header to user[%s]\n", reqId, template.Name, user)
			conn.Channel <- model.UserUpdate{
				Type:   model.ChatUpdate,
				ChatID: template.ID,
				Author: template.ActiveUser,
				Msg:    shortHtml,
			}
		}(user)
	}
	wg.Wait()
	if len(errors) > 0 {
		return fmt.Errorf("%+v", errors)
	}
	return nil
}

func (c *ChatController) InviteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> InviteUser\n", utils.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("--%s-> InviteUser TRACE check login\n", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	chatID, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR chat id, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invitee := r.FormValue("invitee")
	log.Printf("--%s-> InviteUser TRACE inviting[%s] to chat[%d]\n", utils.GetReqId(r), invitee, chatID)
	err = app.State.InviteUser(user, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chat, err := app.State.GetChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR cannot find chat[%d], %s\n", utils.GetReqId(r), chatID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temlate := chat.ToTemplate(user)
	sendChatHeader(utils.GetReqId(r), temlate)

	log.Printf("<-%s-- InviteUser TRACE user [%s] added to chat [%d]\n", utils.GetReqId(r), invitee, chatID)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf(" [%s] ", invitee)))
}

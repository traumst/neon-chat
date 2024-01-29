package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.chat/models"
	"go.chat/utils"
)

var chats = models.ChatList{}

func PollUpdates(w http.ResponseWriter, r *http.Request, user string) {
	lastPing := time.Now()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	utils.SetSseHeaders(w)

	userChats := chats.GetChats(user)
	loopCount := 0
	msgCount := 0
	for {
		loopCount += 1
		log.Printf("=-%s-- PollUpdates TRACE loop %d\n", utils.GetReqId(r), loopCount)
		select {
		case <-r.Context().Done():
			log.Printf("<-%s-- PollUpdates INFO conn closed\n", utils.GetReqId(r))
			return
		case <-ticker.C:
			log.Printf("=-%s-- PollUpdates INFO distribute nsgs\n", utils.GetReqId(r))
			userChats, msgCount = DistributeMessages(w, r, user, userChats)
		}

		if loopCount > stop {
			log.Printf("<-%s-- PollUpdates TRACE stop SSE after %d loops\n", utils.GetReqId(r), loopCount)
			return
		}
		if msgCount > stop {
			log.Printf("<-%s-- PollUpdates TRACE stop SSE after %d msgs\n", utils.GetReqId(r), msgCount)
			return
		}

		if time.Since(lastPing) > 1*time.Second {
			log.Printf("<-%s-- PollUpdates TRACE ping\n", utils.GetReqId(r))
			fmt.Fprintf(w, "event: ping\n")
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
			w.(http.Flusher).Flush()
			lastPing = time.Now()
		}
	}
}

func DistributeMessages(w http.ResponseWriter, r *http.Request, user string, userChats []*models.Chat) ([]*models.Chat, int) {
	log.Printf("--%s-- DistributeMessages TRACE SSE, %s\n", utils.GetReqId(r), user)
	messageCount := 0
	newChats := chats.GetExcept(user, getChatIds(userChats))
	log.Printf("--%s-- DistributeMessages TRACE looping over %d chats for %s\n", utils.GetReqId(r), len(newChats), user)
	for _, chat := range newChats {
		if chat == nil || chat.Name == "" {
			log.Printf("--%s-- DistributeMessages INFO chat name is empty, %s\n", utils.GetReqId(r), chat.Log())
			continue
		}

		html, error := chat.ToTemplate(user).GetShortHTML()
		log.Printf("--%s-- ChatController.DistributeMessages TRACE rendered:[%s]\n", utils.GetReqId(r), html)
		if error != nil {
			log.Printf("--%s-- DistributeMessages ERROR template, %s\n", utils.GetReqId(r), error)
			continue
		}
		log.Printf("--%s-- ChatController.DistributeMessages TRACE send SSE id:%d,event:%s\n",
			utils.GetReqId(r), messageCount, models.ChatEventName)
		fmt.Fprintf(w, "id: %d\n\n", messageCount)
		fmt.Fprintf(w, "event: %s\n", models.ChatEventName)
		// must escape newlines in SSE
		html = strings.ReplaceAll(html, "\n", " ")
		fmt.Fprintf(w, "data: %s\n\n", html)
		w.(http.Flusher).Flush()
		messageCount += 1
	}
	return append(userChats, newChats...), messageCount
}

func getChatIds(chats []*models.Chat) []int {
	chatIDs := make([]int, 0, len(chats))
	for _, chat := range chats {
		chatIDs = append(chatIDs, chat.ID)
	}
	return chatIDs
}

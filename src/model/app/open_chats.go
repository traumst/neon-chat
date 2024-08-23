package app

import (
	"log"
	"sync"
)

type OpenChats struct {
	mu sync.Mutex
	// userId -> []chatId
	open map[uint][]uint
}

func NewOpenChats() *OpenChats {
	cl := OpenChats{
		open: make(map[uint][]uint),
	}
	return &cl
}

func (cl *OpenChats) GetOpenChats(userId uint) []uint {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	return cl.open[userId]
}

func (cl *OpenChats) OpenChat(userId uint, chatId uint) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if cl.open[userId] == nil {
		log.Printf("OpenChat TRACE user[%d] has no open chats, opening chat[%d]", userId, chatId)
		cl.open[userId] = []uint{chatId}
		return
	}
	for _, userChats := range cl.open[userId] {
		if userChats == chatId {
			log.Printf("OpenChat INFO user[%d] already has chat[%d] open", userId, chatId)
			return
		}
	}
	cl.open[userId] = append(cl.open[userId], chatId)
}

func (cl *OpenChats) CloseChat(userId uint, chatId uint) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	openChats := cl.open[userId]
	if len(openChats) <= 0 {
		log.Printf("CloseChat TRACE user[%d] has no open chats", userId)
		return
	}
	for i, openChat := range openChats {
		if openChat == chatId {
			log.Printf("CloseChat INFO user[%d] closes chat[%d]", userId, chatId)
			cl.open[userId] = append(openChats[:i], openChats[i+1:]...)
			return
		}
	}
	log.Printf("CloseChat INFO user[%d] does not have chat[%d]", userId, chatId)
}

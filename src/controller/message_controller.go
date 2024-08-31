package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	h "neon-chat/src/utils/http"
)

func QuoteMessage(
	state *state.State,
	db *d.DBConn,
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("[%s] QuoteMessage TRACE\n", h.GetReqId(r))
	if r.Method != "GET" {
		log.Printf("[%s] QuoteMessage ERROR request method\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	user, _ := handler.ReadSession(state, db, w, r)
	if user == nil {
		log.Printf("[%s] QuoteMessage ERROR user has no session\n", h.GetReqId(r))
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	chatId, msgId, err := parseQuotedMessageArgs(r)
	if err != nil || chatId == 0 || msgId == 0 {
		log.Printf("[%s] QuoteMessage ERROR parsing arguments, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	tmplMsg, err := handler.HandleGetMessage(state, db, user, chatId, msgId)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR quoting message, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to quote message"))
		return
	}
	html, err := tmplMsg.ShortHTML()
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR templating quote, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template message"))
		return
	}
	log.Printf("[%s] QuoteMessage done\n", h.GetReqId(r))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))

	log.Printf("[%s] QuoteMessage TRACE serving html\n", h.GetReqId(r))
	w.WriteHeader(http.StatusOK)
}

func parseQuotedMessageArgs(r *http.Request) (chatId uint, msgId uint, err error) {
	// supports multiple values for the same key, ie delete in bulk, etc
	args := r.URL.Query()
	for k, v := range args {
		switch k {
		case "chatid":
			c, e := strconv.Atoi(v[0])
			if e != nil {
				err = e
			} else {
				chatId = uint(c)
			}
		case "msgid":
			m, e := strconv.Atoi(v[0])
			if e != nil {
				err = e
			} else {
				msgId = uint(m)
			}
		default:
			log.Printf("[%s] WARN parseQuotedMessageArgs unknown argument - [%s:%s]\n", h.GetReqId(r), k, v[0])
		}
		if err != nil {
			log.Printf("[%s] ERROR parseQuotedMessageArgs bad argument - [%s:%s]\n", h.GetReqId(r), k, v[0])
			return 0, 0, fmt.Errorf("invalid argument")
		}
	}
	return chatId, msgId, err
}

func AddMessage(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] TRACE AddMessage \n", h.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("[%s] AddMessage ERROR request method\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	author, err := handler.ReadSession(state, db, w, r)
	if err != nil || author == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] WARN AddMessage bad argument - chatid\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	msg, err := handler.FormValueString(r, "msg")
	if err != nil || len(msg) < 1 {
		log.Printf("[%s] WARN AddMessage bad argument - msg\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}

	message, err := handler.HandleMessageAdd(state, db, chatId, author, msg)
	if err != nil || message == nil {
		log.Printf("[%s] AddMessage ERROR while handing, %s, %v\n", h.GetReqId(r), err.Error(), message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}

	log.Printf("[%s] AddMessage TRACE serving html\n", h.GetReqId(r))
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] DeleteMessage\n", reqId)
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - chatid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := handler.FormValueUint(r, "msgid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - msgid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleted, err := handler.HandleMessageDelete(state, db, chatId, user, msgId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR deletion failed: %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deleted == nil {
		log.Printf("[%s] DeleteMessage WARN message[%d] not found in chat[%d] for user[%d]\n",
			reqId, msgId, chatId, user.Id)
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}

	log.Printf("[%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("~~deleted~~"))
}

package handler

import (
	"fmt"
	"log"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
	"net/http"
	"strconv"
)

type QueryArgs struct {
	ChatId uint
	MsgId  uint
}

func QueryStringArgs(r *http.Request) (parsed QueryArgs, err error) {
	args := r.URL.Query()
	// v is []string, but we only support one value per key
	for k, v := range args {
		switch k {
		case "chatid":
			c, e := strconv.Atoi(v[0])
			if e != nil {
				err = e
			} else {
				parsed.ChatId = uint(c)
			}
		case "msgid":
			m, e := strconv.Atoi(v[0])
			if e != nil {
				err = e
			} else {
				parsed.MsgId = uint(m)
			}
		default:
			log.Printf("[%s] WARN ParseQueryArgs unknown argument - [%s:%s]\n", h.GetReqId(r), k, v[0])
		}
		if err != nil {
			log.Printf("[%s] ERROR ParseQueryArgs bad argument - [%s:%v]\n", h.GetReqId(r), k, v)
			err = fmt.Errorf("invalid argument [%s:%v], %s", k, v, err)
		}
	}
	return parsed, err
}

func FormValueUint(r *http.Request, key string) (uint, error) {
	rawVal := r.PostFormValue(key)
	if rawVal == "" {
		return 0, fmt.Errorf("failed to read key[%s]", key)
	}
	val, err := strconv.Atoi(rawVal)
	if err != nil {
		return 0, fmt.Errorf("failed to parse key[%s] value[%s]", key, rawVal)
	}
	return uint(val), nil
}

func FormValueString(r *http.Request, key string) (string, error) {
	rawVal := r.PostFormValue(key)
	if rawVal == "" {
		return "", fmt.Errorf("failed to read key[%s]", key)
	}
	rawVal = utils.ReplaceWithSingleSpace(rawVal)
	rawVal = utils.RemoveSpecialChars(rawVal)
	return rawVal, nil
}

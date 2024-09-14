package shared

import (
	"fmt"
	"log"
	"neon-chat/src/consts"
	"net/http"
	"strconv"
)

type QueryArgs struct {
	ChatId uint
	MsgId  uint
}

func ParseQueryString(r *http.Request) (parsed QueryArgs, err error) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
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
			log.Printf("[%s] WARN ParseQueryArgs unknown argument - [%s:%s]\n", reqId, k, v[0])
		}
		if err != nil {
			log.Printf("[%s] ERROR ParseQueryArgs bad argument - [%s:%v]\n", reqId, k, v)
			err = fmt.Errorf("invalid argument [%s:%v], %s", k, v, err)
		}
	}
	return parsed, err
}

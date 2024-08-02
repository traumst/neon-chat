package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func FormValueUint(r *http.Request, key string) (uint, error) {
	rawVal := r.PostFormValue(key)
	if rawVal == "" {
		log.Printf("FormValueUint ERROR \n", key)
		return 0, fmt.Errorf("failed to read key[%s]", key)
	}
	parsedVal, err := strconv.Atoi(rawVal)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR parse chatid, %s\n", err.Error())
		return 0, fmt.Errorf("failed to parse key[%s] value[%s]", key, parsedVal)
	}
	return uint(parsedVal), nil
}

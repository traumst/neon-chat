package handler

import (
	"fmt"
	"net/http"
	"prplchat/src/utils"
	"strconv"
)

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
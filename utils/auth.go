package utils

import (
	"fmt"
	"net/http"
	"time"
)

type AuthError struct {
	msg string
}

func (e *AuthError) Error() string {
	return e.msg
}

const usernameCookie = "username"

func AuthUser(w http.ResponseWriter, r *http.Request) (string, error) {
	username := r.FormValue(usernameCookie)
	if username != "" {
		http.SetCookie(w, &http.Cookie{
			Name:    usernameCookie,
			Value:   username,
			Expires: time.Now().Add(8 * time.Hour),
		})
		return username, nil
	}
	return "MISSING_USERNAME", &AuthError{"MISSING_USERNAME"}
}

func GetCurrentUser(r *http.Request) (string, error) {
	username, err := r.Cookie(usernameCookie)
	if err != nil {
		return "ERROR", &AuthError{fmt.Sprintf("UNAUTHORIZED %s", err.Error())}
	}

	return username.Value, nil
}

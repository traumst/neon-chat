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

type Auth struct{}

func AuthUser(w http.ResponseWriter, r *http.Request) (string, error) {
	username := r.FormValue(UsernameCookie)
	if username != "" {
		http.SetCookie(w, &http.Cookie{
			Name:    UsernameCookie,
			Value:   username,
			Expires: time.Now().Add(8 * time.Hour),
		})
		return username, nil
	}
	return "MISSING_USERNAME", &AuthError{"MISSING_USERNAME"}
}

func GetCurrentUser(r *http.Request) (string, error) {
	username, err := r.Cookie(UsernameCookie)
	if err != nil {
		return "ERROR", &AuthError{fmt.Sprintf("UNAUTHORIZED %s", err.Error())}
	}

	return username.Value, nil
}

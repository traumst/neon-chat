package utils

import (
	"fmt"
	"net/http"
)

type AuthError struct {
	msg string
}

func (e *AuthError) Error() string {
	return e.msg
}

type Auth struct{}

func GetCurrentUser(r *http.Request) (string, error) {
	username, err := r.Cookie(SessionCookie)
	if err != nil {
		return "ERROR", &AuthError{fmt.Sprintf("UNAUTHORIZED %s", err.Error())}
	}

	return username.Value, nil
}

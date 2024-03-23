package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.chat/model/app"
)

type Session struct {
	UserId   uint
	UserType app.UserType
	AuthType app.AuthType
}

func (s Session) String() string {
	return fmt.Sprintf("%d:%s:%s:%s", s.UserId, s.UserType, RandStringBytes(9), s.AuthType)
}

func fromString(s string) (*Session, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 4 {
		return nil, fmt.Errorf("invalid session")
	}
	if ss[0] == "" || ss[1] == "" || ss[2] == "" || ss[3] == "" {
		return nil, fmt.Errorf("invalid session content")
	}
	userId, err := strconv.Atoi(ss[0])
	if err != nil {
		return nil, fmt.Errorf("invaild session userId, %s", err)
	}

	return &Session{
		UserId:   uint(userId),
		UserType: app.UserType(ss[1]),
		AuthType: app.AuthType(ss[3]),
	}, nil
}

func GetSessionCookie(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return nil, err
	}

	return fromString(cookie.Value)
}

func SetSessionCookie(w http.ResponseWriter, user *app.User, auth *app.UserAuth, expiration time.Time) {
	cookie := Session{
		UserId:   user.Id,
		UserType: user.Type,
		AuthType: auth.Type,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    SessionCookie,
		Value:   cookie.String(),
		Expires: expiration,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    SessionCookie,
		Value:   "",
		Expires: time.Now(),
	})
}

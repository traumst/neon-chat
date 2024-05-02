package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.chat/src/model/app"
	"go.chat/src/utils"
)

type Session struct {
	UserId   uint
	UserType app.UserType
	AuthType app.AuthType
}

func (s Session) String() string {
	cookie := fmt.Sprintf("%d:%s:%s:%s", s.UserId, s.UserType, utils.RandStringBytes(9), s.AuthType)
	return base64.StdEncoding.EncodeToString([]byte(cookie))
}

func GetSessionCookie(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(utils.SessionCookie)
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
		Name:    utils.SessionCookie,
		Value:   cookie.String(),
		Expires: expiration,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    utils.SessionCookie,
		Value:   "",
		Expires: time.Now(),
	})
}

func fromString(s string) (*Session, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid session, %s", err)
	}
	ss := strings.Split(string(decoded), ":")
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

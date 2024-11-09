package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"neon-chat/src/app"
	"neon-chat/src/app/enum"
	"neon-chat/src/consts"
	"neon-chat/src/utils"
)

var sessions map[uint]Session = map[uint]Session{}

type Session struct {
	UserId   uint
	UserType enum.UserType
	AuthType enum.AuthType
	Expire   time.Time
}

func GetSessionCookie(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		return nil, err
	}
	session, err := decode(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session[%s], %s", cookie.Value, err.Error())
	}
	cached, ok := sessions[session.UserId]
	if !ok {
		return &cached, fmt.Errorf("user has no cached session")
	} else if cached.UserId != session.UserId {
		return &cached, fmt.Errorf("user id mismatch")
	} else if cached.Expire.Before(time.Now()) {
		delete(sessions, session.UserId)
		return &cached, fmt.Errorf("session expired")
	}
	return &cached, nil
}

func SetSessionCookie(w http.ResponseWriter, user *app.User, auth *app.Auth) Session {
	expiration := time.Now().Add(8 * time.Hour)
	session := Session{
		UserId:   user.Id,
		UserType: user.Type,
		AuthType: auth.Type,
		Expire:   expiration,
	}
	sessions[user.Id] = session
	cookie := encode(session)
	http.SetCookie(w, &http.Cookie{
		Name:    consts.SessionCookie,
		Value:   cookie,
		Expires: expiration,
	})
	return session
}

func ClearSessionCookie(w http.ResponseWriter, userId uint) {
	delete(sessions, userId)
	http.SetCookie(w, &http.Cookie{
		Name:    consts.SessionCookie,
		Value:   "",
		Expires: time.Now(),
	})
}

func SaveSessionsToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(sessions)
}

func LoadSessionsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&sessions)
}

func encode(s Session) string {
	cookie := fmt.Sprintf("%d:%s:%s:%s", s.UserId, s.UserType, utils.RandStringBytes(9), s.AuthType)
	return base64.StdEncoding.EncodeToString([]byte(cookie))
}

func decode(s string) (*Session, error) {
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
		UserType: enum.UserType(ss[1]),
		AuthType: enum.AuthType(ss[3]),
	}, nil
}

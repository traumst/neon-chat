package convert

import (
	"neon-chat/src/db"
	"neon-chat/src/model/app"
)

func AuthAppToDB(auth *app.Auth) *db.Auth {
	return &db.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   string(auth.Type),
		Hash:   auth.Hash,
	}
}

func AuthDBToApp(auth *db.Auth) *app.Auth {
	return &app.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   app.AuthType(auth.Type),
		Hash:   auth.Hash,
	}
}

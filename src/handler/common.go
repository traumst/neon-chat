package handler

import (
	"go.chat/src/db"
	"go.chat/src/model/app"
)

func UserToDB(user app.User) db.User {
	return db.User{
		Id:   user.Id,
		Name: user.Name,
		Type: string(user.Type),
		Salt: user.Salt,
	}
}

func UserFromDB(user db.User) app.User {
	return app.User{
		Id:   user.Id,
		Name: user.Name,
		Type: app.UserType(user.Type),
		Salt: user.Salt,
	}
}

func AuthToDB(auth app.UserAuth) db.UserAuth {
	return db.UserAuth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   string(auth.Type),
		Hash:   auth.Hash,
	}
}

func AuthFromDB(auth db.UserAuth) app.UserAuth {
	return app.UserAuth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   app.AuthType(auth.Type),
		Hash:   auth.Hash,
	}
}

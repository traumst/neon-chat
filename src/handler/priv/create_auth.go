package priv

import (
	"fmt"
	"log"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/utils"

	"github.com/jmoiron/sqlx"
)

func CreateAuth(dbConn sqlx.Ext, user *app.User, pass string, authType app.AuthType) (*app.Auth, error) {
	log.Printf("createAuth TRACE IN user[%d] auth[%s]\n", user.Id, authType)
	hash, err := utils.HashPassword(pass, user.Salt)
	if err != nil {
		return nil, fmt.Errorf("error hashing pass, %s", err)
	}
	log.Printf("createAuth TRACE adding user[%d] auth[%s] hash[%s]\n", user.Id, authType, hash)
	dbAuth := &db.Auth{
		Id:     0,
		UserId: user.Id,
		Type:   string(authType),
		Hash:   hash,
	}
	dbAuth, err = db.AddAuth(dbConn, *dbAuth)
	if err != nil || dbAuth == nil {
		return nil, fmt.Errorf("fail to add auth to user[%d][%s], %s", user.Id, user.Name, err)
	}
	if dbAuth.Id <= 0 {
		return nil, fmt.Errorf("user[%d][%s] auth was not created", user.Id, user.Name)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appAuth, err
}

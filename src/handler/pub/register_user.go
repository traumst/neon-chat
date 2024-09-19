package pub

import (
	"fmt"
	"log"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/utils"

	"github.com/jmoiron/sqlx"
)

func RegisterUser(
	dbConn *db.DBConn,
	newUser *app.User,
	pass string,
	authType app.AuthType,
) (*app.User, *app.Auth, error) {
	log.Printf("Register TRACE IN user\n")
	if dbConn == nil || newUser == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", newUser, dbConn)
	}
	if newUser.Name == "" || pass == "" || newUser.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] salt[%s]", newUser.Name, newUser.Salt)
	}
	var appUser *app.User
	var err error
	if newUser.Id != 0 {
		appUser = newUser
	} else {
		appUser, err = createUser(dbConn.Tx, newUser)
		if err != nil || appUser == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", appUser.Name)
		}
	}
	auth, err := createAuth(dbConn.Tx, appUser, pass, authType)
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, appUser, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", appUser.Id, auth)
	return appUser, auth, nil
}

func createAuth(dbConn sqlx.Ext, user *app.User, pass string, authType app.AuthType) (*app.Auth, error) {
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

func createUser(dbConn sqlx.Ext, user *app.User) (*app.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("createUser TRACE creating user[%s]\n", user.Name)
	dbUser := convert.UserAppToDB(user)
	created, err := db.AddUser(dbConn, dbUser)
	if err != nil || created == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", created, err)
	}
	if created.Id <= 0 {
		return nil, fmt.Errorf("user[%s] was not created", created.Name)
	}
	appUser := convert.UserDBToApp(created, nil)
	return appUser, err
}

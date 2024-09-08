package shared

import (
	"fmt"
	"log"

	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func CreateUser(db *d.DBConn, user *a.User) (*a.User, error) {
	if user.Id != 0 && user.Salt != "" {
		log.Printf("createUser TRACE completing user[%s] signup\n", user.Name)
		return user, nil
	}
	log.Printf("createUser TRACE creating user[%s]\n", user.Name)
	dbUser := convert.UserAppToDB(user)
	created, err := db.AddUser(dbUser)
	if err != nil || created == nil {
		return nil, fmt.Errorf("failed to add user[%v], %s", created, err)
	}
	if created.Id <= 0 {
		return nil, fmt.Errorf("user[%s] was not created", created.Name)
	}
	appUser := convert.UserDBToApp(created, nil)
	return appUser, err
}

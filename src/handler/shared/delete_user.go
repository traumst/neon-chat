package shared

import (
	"fmt"
	"log"

	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func DeleteUser(db *d.DBConn, user *a.User) error {
	if user.Id < 1 {
		log.Printf("deleteUser TRACE completing user[%s] signup\n", user.Name)
		return nil
	}
	log.Printf("deleteUser TRACE creating user[%s]\n", user.Name)
	err := db.DropUser(user.Id)
	if err != nil {
		return fmt.Errorf("failed to delete, %s", err)
	}
	return nil
}

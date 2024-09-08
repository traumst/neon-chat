package shared

import (
	"fmt"
	"log"
	d "neon-chat/src/db"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
)

func UpdateUser(state *state.State, db *d.DBConn, payload *a.User) (*a.User, error) {
	current, err := GetUser(db, payload.Id)
	if err != nil {
		return nil, fmt.Errorf("user for update[%d] not found: %s", payload.Id, err.Error())
	}

	if current.Status != payload.Status {
		current.Status = payload.Status
		err = db.UpdateUserStatus(payload.Id, string(payload.Status))
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] status: %s", payload.Id, err.Error())
		}
	}
	if current.Name != payload.Name {
		current.Name = payload.Name
		err = db.UpdateUserName(payload.Id, payload.Name)
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] name: %s", payload.Id, err.Error())
		}
	}

	return current, nil
}

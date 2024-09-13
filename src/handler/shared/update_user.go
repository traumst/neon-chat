package shared

import (
	"fmt"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/state"
	"neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func UpdateUser(state *state.State, dbConn sqlx.Ext, payload *app.User) (*app.User, error) {
	current, err := GetUser(dbConn, payload.Id)
	if err != nil {
		return nil, fmt.Errorf("user for update[%d] not found: %s", payload.Id, err.Error())
	}

	if current.Status != payload.Status {
		current.Status = payload.Status
		err = db.UpdateUserStatus(dbConn, payload.Id, string(payload.Status))
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] status: %s", payload.Id, err.Error())
		}
	}
	if current.Name != payload.Name {
		current.Name = payload.Name
		err = db.UpdateUserName(dbConn, payload.Id, payload.Name)
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] name: %s", payload.Id, err.Error())
		}
	}

	return current, nil
}

package handler

import (
	"fmt"
	"log"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
)

// empty targetUser means all users in chat
func DistributeUserChange(
	state *AppState,
	//TODO targetUser *app.User, // who to inform, nil for all users
	subjectUser *app.User, // which user changed, nil for every user in chat
	updateType event.UpdateType,
) error {
	if subjectUser == nil {
		return fmt.Errorf("subject user is nil")
	}
	// TODO distribute to users with common chats
	targetUser := subjectUser
	err := distributeUpdateOfUser(
		state,
		targetUser,
		subjectUser,
		updateType)
	return err
}

func distributeUpdateOfUser(
	state *AppState,
	targetUser *app.User,
	subjectUser *app.User,
	updateType event.UpdateType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked: %v", r)
		}
	}()

	var conn *Conn
	conn, err = state.GetConn(targetUser.Id)
	if err != nil {
		return err
	}
	if conn.User.Id != targetUser.Id {
		return fmt.Errorf("user[%d] does not own conn[%v], user[%d] does", targetUser.Id, conn.Origin, conn.User.Id)
	}
	switch updateType {
	case event.UserChange:
		return userNameChanged(conn, subjectUser)
	default:
		return fmt.Errorf("unknown event type[%v]", updateType)
	}
}

func userNameChanged(conn *Conn, subject *app.User) error {
	if subject == nil {
		return fmt.Errorf("subjectUser is nil for userChanged")
	}
	log.Printf("∞----> userChanged TRACE informing target[%d] about subject[%d] change\n", conn.User.Id, subject.Id)
	tmpl := subject.Template(-64, 0, conn.User.Id)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template user")
	}
	conn.In <- event.LiveUpdate{
		Event:    event.UserChange,
		ChatId:   -2,
		UserId:   subject.Id,
		MsgId:    -3,
		AuthorId: subject.Id,
		Data:     data,
	}
	return nil
}

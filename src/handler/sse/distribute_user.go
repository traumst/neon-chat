package sse

import (
	"fmt"
	"log"

	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/event"
)

// empty targetUser means all users in chat
func DistributeUserChange(
	state *state.State,
	//TODO targetUser *app.User, // who to inform, nil for all users
	subjectUser *app.User, // which user changed, nil for every user in chat
	updateType event.EventType,
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
	state *state.State,
	targetUser *app.User,
	subjectUser *app.User,
	updateType event.EventType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked: %v", r)
		}
	}()
	conns := state.GetConn(targetUser.Id)
	for _, conn := range conns {
		if conn.User.Id != targetUser.Id {
			return fmt.Errorf("user[%d] does not own conn[%v], user[%d] does", targetUser.Id, conn.Origin, conn.User.Id)
		}
		var connerr error
		switch updateType {
		case event.UserChange:
			connerr = userNameChanged(conn, subjectUser)
		default:
			connerr = fmt.Errorf("unknown event type[%v]", updateType)
		}
		if err == nil {
			err = connerr
		} else {
			err = fmt.Errorf("%s, %s", err.Error(), connerr.Error())
		}
	}
	return err
}

func userNameChanged(conn *state.Conn, subject *app.User) error {
	if subject == nil {
		return fmt.Errorf("subjectUser is nil for userChanged")
	}
	log.Printf("userChanged TRACE informing target[%d] about subject[%d] change\n", conn.User.Id, subject.Id)
	tmpl := subject.Template(-64, 0, conn.User.Id)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template user")
	}
	conn.In <- event.LiveEvent{
		Event:    event.UserChange,
		ChatId:   -2,
		UserId:   subject.Id,
		MsgId:    -3,
		AuthorId: subject.Id,
		Data:     data,
	}
	return nil
}

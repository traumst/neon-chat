package handler

import (
	"fmt"
	"log"

	"prplchat/src/model/app"
	"prplchat/src/model/event"
)

func DistributeAvatarChange(
	state *AppState,
	//TODO targetUser *app.User, // who to inform, nil for all users
	subjectUser *app.User, // which user changed, nil for every user in chat
	avatar *app.Avatar,
	updateType event.UpdateType,
) error {
	if subjectUser == nil {
		return fmt.Errorf("subject user is nil")
	}
	// TODO distribute to users with common chats
	targetUser := subjectUser
	err := distributeUpdateOfAvatar(
		state,
		targetUser,
		subjectUser,
		avatar,
		updateType)
	return err
}

func distributeUpdateOfAvatar(
	state *AppState,
	targetUser *app.User,
	subjectUser *app.User,
	avatar *app.Avatar,
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
	case event.AvatarChange:
		return avatarChanged(conn, subjectUser, avatar)
	default:
		return fmt.Errorf("unknown event type[%v]", updateType)
	}
}

func avatarChanged(conn *Conn, subject *app.User, avatar *app.Avatar) error {
	if subject == nil || avatar == nil {
		return fmt.Errorf("arguments were nil, user[%v], avatar[%v]", subject, avatar)
	}
	log.Printf("avatarChanged TRACE informing target[%d] about subject[%d] new avatar[%d]\n",
		conn.User.Id, subject.Id, avatar.Id)
	tmpl := avatar.Template(subject)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template user")
	}
	conn.In <- event.LiveUpdate{
		Event:    event.AvatarChange,
		ChatId:   -2,
		UserId:   subject.Id,
		MsgId:    -3,
		AuthorId: subject.Id,
		Data:     data,
	}
	return nil
}

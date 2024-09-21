package sse

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/event"
	"neon-chat/src/state"
)

func DistributeAvatarChange(
	state *state.State,
	user *app.User, // which user changed
	avatar *app.Avatar,
	updateType event.EventType,
) error {
	if user == nil {
		return fmt.Errorf("subject user is nil")
	}
	err := distributeUpdateOfAvatar(
		state,
		user,
		avatar,
		updateType)
	return err
}

func distributeUpdateOfAvatar(
	state *state.State,
	user *app.User,
	avatar *app.Avatar,
	updateType event.EventType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked: %v", r)
		}
	}()
	conns := state.GetConn(user.Id)
	if len(conns) == 0 {
		return nil
	}
	switch updateType {
	case event.AvatarChange:
		return avatarChanged(conns, user, avatar)
	default:
		return fmt.Errorf("unknown event type[%v]", updateType)
	}
}

func avatarChanged(conns []*state.Conn, user *app.User, avatar *app.Avatar) error {
	if len(conns) == 0 {
		return nil
	}
	if user == nil || avatar == nil {
		return fmt.Errorf("arguments were nil, user[%v], avatar[%v]", user, avatar)
	}
	log.Printf("avatarChanged TRACE informing target[%d] about subject[%d] new avatar[%d]\n",
		conns[0].User.Id, user.Id, avatar.Id)
	tmpl := avatar.Template(user)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template user")
	}
	for _, conn := range conns {
		conn.In <- event.LiveEvent{
			Event:    event.AvatarChange,
			ChatId:   0,
			UserId:   user.Id,
			MsgId:    0,
			AuthorId: user.Id,
			Data:     data,
		}
	}
	return nil
}

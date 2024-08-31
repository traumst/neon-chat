package sse

import (
	"fmt"
	"log"
	"sync"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/state"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
)

// targetUser=nil means all users in chat
func DistributeUserChange(
	state *state.State,
	db *db.DBConn,
	targetUser *app.User, // who to inform, nil for all users
	subjectUser *app.User, // which user changed
	updateType event.EventType,
) error {
	if subjectUser == nil {
		return fmt.Errorf("subject user is nil")
	}
	if targetUser != nil {
		return distributeUpdateOfUser(state, targetUser, subjectUser, updateType)
	}
	dbChats, err := db.GetUserChats(subjectUser.Id)
	if err != nil {
		return fmt.Errorf("failed to get chats for user[%d], %s", subjectUser.Id, err)
	}
	appChats, err := convertToApp(db, dbChats)
	if len(appChats) <= 0 {
		log.Printf("userChanged WARN user[%d] has no chats\n", subjectUser.Id)
		return fmt.Errorf("user[%d] has no chats", subjectUser.Id)
	}
	var wg sync.WaitGroup
	for _, appChat := range appChats {
		if appChat == nil {
			continue
		}
		wg.Add(1)
		go func(chat *app.Chat) {
			defer wg.Done()
			err := distributeInCommonChat(db, chat, state, subjectUser, updateType)
			if err != nil {
				log.Printf("userChanged ERROR failed to distribute to chat[%d], %s\n", chat.Id, err)
			}
		}(appChat)
	}
	wg.Wait()
	return err
}

func convertToApp(db *db.DBConn, dbChats []db.Chat) ([]*app.Chat, error) {
	// chatId by userId
	chatIdToOwnerId := make(map[uint]uint)
	// owners by userId
	ownerIdSet := make(map[uint]bool)
	ownerIds := make([]uint, 0)
	for _, dbChat := range dbChats {
		chatIdToOwnerId[dbChat.Id] = dbChat.OwnerId
		if _, ok := ownerIdSet[dbChat.OwnerId]; !ok {
			ownerIdSet[dbChat.OwnerId] = true
			ownerIds = append(ownerIds, dbChat.OwnerId)
		}
	}
	dbChatOwners, err := db.GetUsers(ownerIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat owners[%v], %s", ownerIds, err)
	}
	ownerMap := make(map[uint]*app.User)
	for _, dbChatOwner := range dbChatOwners {
		ownerMap[dbChatOwner.Id] = convert.UserDBToApp(&dbChatOwner)
	}
	var chats []*app.Chat
	for _, dbChat := range dbChats {
		owner, ok := ownerMap[dbChat.OwnerId]
		if !ok {
			log.Printf("userChanged WARN chat[%d] skipped due to owner[%d] not found\n", dbChat.Id, dbChat.OwnerId)
			continue
		}
		chats = append(chats, convert.ChatDBToApp(&dbChat, owner))
	}
	return chats, nil
}

func distributeInCommonChat(
	db *db.DBConn,
	chat *app.Chat,
	state *state.State,
	subjectUser *app.User,
	updateType event.EventType,
) error {
	dbUsers, err := db.GetChatUsers(chat.Id)
	if err != nil {
		return fmt.Errorf("failed to get users in chat[%d], %s", chat.Id, err)
	}
	var targetUsers []*app.User
	for _, dbUser := range dbUsers {
		targetUsers = append(targetUsers, convert.UserDBToApp(&dbUser))
	}
	// inform if chat is open
	var errs []string
	for _, targetUser := range targetUsers {
		if targetUser == nil {
			continue
		}
		openChatId := state.GetOpenChat(targetUser.Id)
		if openChatId == 0 || openChatId != chat.Id {
			// TODO send unread counter update
			continue
		}
		uErr := distributeUpdateOfUser(state, targetUser, subjectUser, updateType)
		if uErr != nil {
			errs = append(errs, fmt.Sprintf("failed to distribute to user[%d], %s\n", targetUser.Id, uErr))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to distribute in chat[%d], %v", chat.Id, errs)
	}
	return nil
}

func distributeUpdateOfUser(
	state *state.State,
	targetUser *app.User,
	subjectUser *app.User,
	updateType event.EventType,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DUOU FATAL: %v", r)
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
		return fmt.Errorf("subject user is nil")
	}
	log.Printf("userNameChanged TRACE informing target[%d] about subject[%d] name change\n", conn.User.Id, subject.Id)
	tmpl := subject.Template(0, 0, conn.User.Id)
	data, err := tmpl.HTML()
	if err != nil {
		return fmt.Errorf("failed to template subject user")
	}
	conn.In <- event.LiveEvent{
		Event:    event.UserChange,
		ChatId:   0,
		UserId:   subject.Id,
		MsgId:    0,
		AuthorId: subject.Id,
		Data:     data,
	}
	return nil
}

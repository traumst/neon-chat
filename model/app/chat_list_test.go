package app

import "testing"

func TestInit(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	if cl.isInit {
		t.Errorf("TestAddChat expected isInit false, got true")
		return
	}
	cl.init()
	if !cl.isInit {
		t.Errorf("TestAddChat expected isInit true, got false")
	}
}

func TestAddChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID1 := cl.AddChat(&user, "test-chat")
	if chatID1 != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID1)
		return
	}
	user2 := User{Id: 2, Name: "Jill", Type: UserType(Free)}
	chatID2 := cl.AddChat(&user2, "test-chat")
	if chatID2 != 1 {
		t.Errorf("TestAddChat expected chatID 1, got %d", chatID2)
	} else if chatID2 == chatID1 {
		t.Errorf("TestAddChat added chat with duplicate id %d", chatID2)
	}
}

func TestOpenChat(t *testing.T) {
	t.Logf("TestOpenChat started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID1 := cl.AddChat(&user, "test-chat")

	user2 := User{Id: 2, Name: "Jill", Type: UserType(Free)}
	chatID2 := cl.AddChat(&user2, "test-chat-2")

	if chatID2 == chatID1 {
		t.Errorf("TestAddChat added chat with duplicate id %d", chatID2)
		return
	}
	openChat2 := cl.GetOpenChat(user2.Id)
	if openChat2 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	openChat1, err := cl.OpenChat(user.Id, chatID1)
	if err != nil {
		t.Errorf("TestOpenChat failed to open chat [%s]", err)
		return
	}
	if openChat1 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	if openChat1.Id == openChat2.Id || openChat1.Id != chatID1 {
		t.Errorf("TestOpenChat expected chatID %d, got %d", chatID1, openChat1.Id)
	}
}

func TestGetChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID := cl.AddChat(&user, "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID)
		return
	}
	chats := cl.GetChats(user.Id)
	if chats == nil {
		t.Errorf("TestAddChat chats were nil")
	} else if len(chats) != 1 {
		t.Errorf("TestAddChat expected 1 chat, got %d", len(chats))
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	t.Logf("TestGetOpenChatEmpty started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chat := cl.GetOpenChat(user.Id)
	if chat != nil {
		t.Errorf("TestGetOpenChatEmpty expected chat to be NIL")
	}
}

func TestGetOpenChat(t *testing.T) {
	t.Logf("TestGetOpenChat started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID := cl.AddChat(&user, "test-chat")
	if chatID != 0 {
		t.Errorf("TestGetOpenChat expected chatID 0, got [%d]", chatID)
		return
	}
	chat := cl.GetOpenChat(user.Id)
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
		return
	}
	chatID = cl.AddChat(&user, "test-chat-2")
	if chatID != 1 {
		t.Errorf("TestGetOpenChat expected chatID 1, got [%d]", chatID)
		return
	}
	chat = cl.GetOpenChat(user.Id)
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
	} else if chat.Id != chatID {
		t.Errorf("TestGetOpenChat expected chatID 1, got [%d]", chat.Id)
	}
}

func TestGetChats(t *testing.T) {
	t.Logf("TestGetChats started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	_ = cl.AddChat(&user, "test-chat")
	_ = cl.AddChat(&user, "test-chat2")
	_ = cl.AddChat(&user, "test-chat3")
	chats := cl.GetChats(user.Id)
	if chats == nil {
		t.Errorf("TestGetChats chats were nil")
	} else if len(chats) != 3 {
		t.Errorf("TestGetChats expected 1 chat, got [%d]", len(chats))
	}
}

func TestDeleteChatEmpty(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	err := cl.DeleteChat(user.Id, nil)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
	}
	chat1 := Chat{Id: 1, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat1)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 1")
	}
	chat2 := Chat{Id: 0, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat2)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 0")
	}
	chat3 := Chat{Id: -1, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat3)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, -1")
	}
}

func TestDeleteChat(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID := cl.AddChat(&user, "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	openChat := cl.GetOpenChat(user.Id)
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat(user.Id, chatID)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat(user.Id, chat)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty delete failed [%s]", err)
		return
	}
	openChat = cl.GetOpenChat(user.Id)
	if openChat != nil {
		t.Errorf("TestDeleteChatEmpty openChat was expected to be NIL, but was [%+v]", openChat)
	}
}

func TestDeleteChatNotOwner(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	user := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID := cl.AddChat(&user, "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	openChat := cl.GetOpenChat(user.Id)
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat(user.Id, chatID)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat(333, chat)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty delete should not have been allowed")
	}
}

func TestInviteUser(t *testing.T) {
	t.Logf("TestInviteUserEmpty started")
	cl := ChatList{}
	owner := User{Id: 1, Name: "John", Type: UserType(Free)}
	chatID := cl.AddChat(&owner, "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID)
		return
	}
	invitee := User{Id: 2, Name: "Jill", Type: UserType(Free)}
	err := cl.InviteUser(owner.Id, chatID, &invitee)
	if err != nil {
		t.Errorf("TestInviteUserEmpty fail to invite user [%v]", invitee)
		return
	}
	inviteeChats := cl.GetChats(invitee.Id)
	if inviteeChats == nil {
		t.Errorf("TestInviteUserEmpty inviteeChats were nil")
	} else if len(inviteeChats) != 1 {
		t.Errorf("TestInviteUserEmpty expected 1 chat, got %d", len(inviteeChats))
	} else if inviteeChats[0].Id != chatID {
		t.Errorf("TestInviteUserEmpty expected chatID %d, got %d", chatID, inviteeChats[0].Id)
	} else if !inviteeChats[0].isUserInChat(invitee.Id) {
		t.Errorf("TestInviteUserEmpty user [%v] is not in chat", invitee)
	} else if !inviteeChats[0].isUserInChat(owner.Id) {
		t.Errorf("TestInviteUserEmpty owner [%v] is not in chat", owner)
	}
}

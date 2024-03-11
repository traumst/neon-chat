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
	chatID1 := cl.AddChat("test-user", "test-chat")
	if chatID1 != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID1)
		return
	}
	chatID2 := cl.AddChat("test-user-2", "test-chat")
	if chatID2 != 1 {
		t.Errorf("TestAddChat expected chatID 1, got %d", chatID2)
	} else if chatID2 == chatID1 {
		t.Errorf("TestAddChat added chat with duplicate id %d", chatID2)
	}
}

func TestOpenChat(t *testing.T) {
	t.Logf("TestOpenChat started")
	cl := ChatList{}
	chatID1 := cl.AddChat("test-user", "test-chat")
	chatID2 := cl.AddChat("test-user", "test-chat-2")
	if chatID2 == chatID1 {
		t.Errorf("TestAddChat added chat with duplicate id %d", chatID2)
		return
	}
	openChat2 := cl.GetOpenChat("test-user")
	if openChat2 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	openChat1, err := cl.OpenChat("test-user", chatID1)
	if err != nil {
		t.Errorf("TestOpenChat failed to open chat [%s]", err)
		return
	}
	if openChat1 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	if openChat1.ID == openChat2.ID || openChat1.ID != chatID1 {
		t.Errorf("TestOpenChat expected chatID %d, got %d", chatID1, openChat1.ID)
	}
}

func TestGetChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID)
		return
	}
	chats := cl.GetChats("test-user")
	if chats == nil {
		t.Errorf("TestAddChat chats were nil")
	} else if len(chats) != 1 {
		t.Errorf("TestAddChat expected 1 chat, got %d", len(chats))
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	t.Logf("TestGetOpenChatEmpty started")
	cl := ChatList{}
	chat := cl.GetOpenChat("test-user")
	if chat != nil {
		t.Errorf("TestGetOpenChatEmpty expected chat to be NIL")
	}
}

func TestGetOpenChat(t *testing.T) {
	t.Logf("TestGetOpenChat started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestGetOpenChat expected chatID 0, got [%d]", chatID)
		return
	}
	chat := cl.GetOpenChat("test-user")
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
		return
	}
	chatID = cl.AddChat("test-user", "test-chat-2")
	if chatID != 1 {
		t.Errorf("TestGetOpenChat expected chatID 1, got [%d]", chatID)
		return
	}
	chat = cl.GetOpenChat("test-user")
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
	} else if chat.ID != chatID {
		t.Errorf("TestGetOpenChat expected chatID 1, got [%d]", chat.ID)
	}
}

func TestGetChats(t *testing.T) {
	t.Logf("TestGetChats started")
	cl := ChatList{}
	_ = cl.AddChat("test-user", "test-chat")
	_ = cl.AddChat("test-user", "test-chat2")
	_ = cl.AddChat("test-user", "test-chat3")
	chats := cl.GetChats("test-user")
	if chats == nil {
		t.Errorf("TestGetChats chats were nil")
	} else if len(chats) != 3 {
		t.Errorf("TestGetChats expected 1 chat, got [%d]", len(chats))
	}
}

func TestDeleteChatEmpty(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	err := cl.DeleteChat("test-user", nil)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
	}
	chat1 := Chat{ID: 1}
	err = cl.DeleteChat("test-user", &chat1)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 1")
	}
	chat2 := Chat{ID: 0}
	err = cl.DeleteChat("test-user", &chat2)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 0")
	}
	chat3 := Chat{ID: -1}
	err = cl.DeleteChat("test-user", &chat3)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, -1")
	}
}

func TestDeleteChat(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	openChat := cl.GetOpenChat("test-user")
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat("test-user", chatID)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat("test-user", chat)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty delete failed [%s]", err)
		return
	}
	openChat = cl.GetOpenChat("test-user")
	if openChat != nil {
		t.Errorf("TestDeleteChatEmpty openChat was expected to be NIL, but was [%+v]", openChat)
	}
}

func TestDeleteChatNotOwner(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	openChat := cl.GetOpenChat("test-user")
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat("test-user", chatID)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat("OTHER-user", chat)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty delete should not have been allowed")
	}
}

func TestInviteUser(t *testing.T) {
	t.Logf("TestInviteUserEmpty started")
	cl := ChatList{}
	owner := "test-user"
	chatID := cl.AddChat(owner, "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID)
		return
	}
	invitee := "poopy"
	err := cl.InviteUser(owner, chatID, invitee)
	if err != nil {
		t.Errorf("TestInviteUserEmpty fail to invite user [%s]", invitee)
		return
	}
	inviteeChats := cl.GetChats(invitee)
	if inviteeChats == nil {
		t.Errorf("TestInviteUserEmpty inviteeChats were nil")
	} else if len(inviteeChats) != 1 {
		t.Errorf("TestInviteUserEmpty expected 1 chat, got %d", len(inviteeChats))
	} else if inviteeChats[0].ID != chatID {
		t.Errorf("TestInviteUserEmpty expected chatID %d, got %d", chatID, inviteeChats[0].ID)
	} else if !inviteeChats[0].isUserInChat(invitee) {
		t.Errorf("TestInviteUserEmpty user [%s] is not in chat", invitee)
	} else if !inviteeChats[0].isUserInChat(owner) {
		t.Errorf("TestInviteUserEmpty owner [%s] is not in chat", owner)
	}
}

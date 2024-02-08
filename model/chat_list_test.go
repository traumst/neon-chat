package model

import "testing"

func TestInInit(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	if cl.isInit {
		t.Errorf("TestAddChat expected isInit false, got true")
		return
	}
	cl.init("test-user")
	if !cl.isInit {
		t.Errorf("TestAddChat expected isInit true, got false")
	}
}

func TestAddChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got %d", chatID)
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
	t.Logf("TestGetOpenChatEmpty started")
	cl := ChatList{}
	chatID := cl.AddChat("test-user", "test-chat")
	if chatID != 0 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	chat := cl.GetOpenChat("test-user")
	if chat == nil {
		t.Errorf("TestGetOpenChatEmpty chat was NIL [%s]", chat.Log())
		return
	}
	chatID = cl.AddChat("test-user", "test-chat-2")
	if chatID != 1 {
		t.Errorf("TestAddChat expected chatID 0, got [%d]", chatID)
		return
	}
	chat = cl.GetOpenChat("test-user")
	if chat == nil {
		t.Errorf("TestGetOpenChatEmpty chat was NIL [%s]", chat.Log())
	} else if chat.ID != chatID {
		t.Errorf("TestGetOpenChatEmpty expected chatID 1, got [%d]", chat.ID)
	}
}

func TestDeleteChatEmpty(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := ChatList{}
	err := cl.DeleteChat("test-user", -1)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
	}
	err = cl.DeleteChat("test-user", 0)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
	}
	err = cl.DeleteChat("test-user", 1)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
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
	err := cl.DeleteChat("test-user", chatID)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty delete failed [%s]", err)
		return
	}
	openChat = cl.GetOpenChat("test-user")
	if openChat != nil {
		t.Errorf("TestDeleteChatEmpty openChat was expected to be NIL, but was [%s]", openChat.Log())
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
	err := cl.DeleteChat("OTHER-user", chatID)
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

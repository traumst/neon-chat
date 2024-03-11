package app

import "testing"

func TestDefaultChat(t *testing.T) {
	t.Logf("TestDefaultChat started")
	c := Chat{}
	if c.ID != 0 {
		t.Errorf("TestDefaultChat expected UnknownUpdate, got [%d]", c.ID)
	} else if c.Name != "" {
		t.Errorf("TestDefaultChat expected empty user, got [%s]", c.Name)
	} else if c.Owner != "" {
		t.Errorf("TestDefaultChat expected empty owner, got [%s]", c.Owner)
	} else if len(c.users) != 0 {
		t.Errorf("TestDefaultChat expected empty users, got [%+v]", c.users)
	}
	t.Logf("TestDefaultChat finished")
}

func TestIsUserInChat(t *testing.T) {
	t.Logf("TestIsUserInChat started")
	u := "test-user"
	o := "other-user"
	c := Chat{users: []string{"user1", u, "user2"}}
	if c.isUserInChat(o) {
		t.Errorf("TestIsUserInChat [%s] should not have been in the chat with [%s]", o, u)
	}
	if !c.isUserInChat(u) {
		t.Errorf("TestIsUserInChat [%s] should have been in the chat", u)
	}
	t.Logf("TestIsUserInChat finished")
}

func TestIsAuthor(t *testing.T) {
	t.Logf("TestIsAuthor started")
	u := "test-user"
	o := "other-user"
	c := Chat{Owner: "owner", users: []string{"user1", u, "user2"}}
	m := Message{ID: 1, Author: u, Text: "test message"}
	mwid, err := c.history.Add(&m)
	if err != nil {
		t.Errorf("TestIsAuthor failed to add message: %s", err)
	}
	if c.isAuthor(o, mwid.ID) {
		t.Errorf("TestIsAuthor [%s] should not have been the author of the message", o)
	}
	if !c.isAuthor(u, mwid.ID) {
		t.Errorf("TestIsAuthor [%s] should have been the author of the message", u)
	}
	t.Logf("TestIsOwner finished")
}

func TestIsOwner(t *testing.T) {
	t.Logf("TestIsOwner started")
	o := "test-owner"
	c := Chat{Owner: o}
	if c.isOwner("other-user") {
		t.Errorf("TestIsOwner [%s] should have been the owner", o)
	}
	if !c.isOwner(o) {
		t.Errorf("TestIsOwner [%s] should have been the owner", o)
	}
	t.Logf("TestIsOwner finished")
}

package model

import (
	"net/http/httptest"
	"testing"

	"go.chat/utils"
)

func NewAppState() *AppState {
	return &AppState{
		chats:    ChatList{},
		userConn: make(UserConn, 0),
	}
}

func TestAddConn(t *testing.T) {
	t.Logf("TestAddConn started")
	state := NewAppState()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := "John"

	conn1 := state.ReplaceConn(w, *r, user)
	if conn1 == nil {
		t.Errorf("Expected a conn1, got nil")
	}

	conn2 := state.ReplaceConn(w, *r, user)
	if conn2 == nil {
		t.Errorf("Expected a conn2, got nil")
	}

	if conn1 == conn2 {
		t.Errorf("Expected conn1 and conn2 to be different")
	}
}

func TestGetConn(t *testing.T) {
	t.Logf("TestGetConn started")
	state := NewAppState()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	utils.SetReqId(r, &reqId)
	user := "John"

	conn := state.ReplaceConn(w, *r, user)
	if conn == nil {
		t.Errorf("TestGetConn expected a conn, got nil")
	}

	conn2, err := state.GetConn(user)
	if err != nil {
		t.Errorf("TestGetConn expected no error, got [%s]", err)
	} else if conn2 == nil {
		t.Errorf("TestGetConn expected a conn2, got nil")
	} else if conn.User != conn2.User ||
		conn.Origin != conn2.Origin ||
		conn.In != conn2.In ||
		conn.Writer != conn2.Writer {
		t.Errorf("TestGetConn expected equality, got [%+v], [%+v]", conn, conn2)
	}
}

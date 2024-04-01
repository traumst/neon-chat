package model

import (
	"net/http/httptest"
	"testing"

	"go.chat/model/app"
	"go.chat/utils"
)

var app1 = &ApplicationState

func TestAddConn(t *testing.T) {
	t.Logf("TestAddConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeFree)}

	conn1 := app1.ReplaceConn(w, *r, &user)
	if conn1 == nil {
		t.Errorf("Expected a conn1, got nil")
	}

	conn2 := app1.ReplaceConn(w, *r, &user)
	if conn2 == nil {
		t.Errorf("Expected a conn2, got nil")
	}

	if conn1 == conn2 {
		t.Errorf("Expected conn1 and conn2 to be different")
	}
}

func TestGetConn(t *testing.T) {
	t.Logf("TestGetConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	utils.SetReqId(r, &reqId)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeFree)}

	conn := app1.ReplaceConn(w, *r, &user)
	if conn == nil {
		t.Errorf("TestGetConn expected a conn, got nil")
	}

	conn2, err := app1.GetConn(user.Id)
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

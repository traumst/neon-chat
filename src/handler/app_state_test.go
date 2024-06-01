package handler

import (
	"math/rand"
	"net/http/httptest"
	"testing"

	"prplchat/src/model/app"
	h "prplchat/src/utils/http"
)

var app1 = &ApplicationState

func TestAddConn(t *testing.T) {
	t.Logf("TestAddConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeLocal)}

	conn1 := app1.AddConn(w, *r, &user)
	if conn1 == nil {
		t.Fatalf("Expected a conn1, got nil")
	}

	conn2 := app1.AddConn(w, *r, &user)
	if conn2 == nil {
		t.Fatalf("Expected a conn2, got nil")
	}

	if conn1 == conn2 {
		t.Fatalf("Expected conn1 and conn2 to be different")
	}
}

func TestGetConn(t *testing.T) {
	t.Logf("TestGetConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	reqId = h.SetReqId(r, &reqId)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeLocal)}
	expect := app1.AddConn(w, *r, &user)
	if expect == nil {
		t.Fatalf("TestGetConn expected a conn, got nil")
	}
	conns := app1.GetConn(user.Id)
	if len(conns) == 0 {
		t.Fatalf("TestGetConn expected connections were empty, count [%d]", len(conns))
	}
	var conn *Conn
	for _, c := range conns {
		if c.Id == expect.Id {
			conn = c
			break
		}
	}

	if expect == nil || conn == nil {
		t.Fatalf("TestGetConn expected connections, got\n expected[%v]\n conn[%v]", expect, conn)
	}
	if expect.User != conn.User {
		t.Fatalf("TestGetConn expected equality,\nexpected user[%v],\nconn origin[%v]", expect.User, conn.User)
	}
	if expect.Origin != conn.Origin {
		t.Fatalf("TestGetConn expected equality,\nexpected origin[%v],\nconn origin[%v]", expect.Origin, conn.Origin)
	}
	if expect.In != conn.In {
		t.Fatalf("TestGetConn expected equality,\nexpected in[%v],\nconn in[%v]", expect.In, conn.In)
	}
	if expect.Writer != conn.Writer {
		t.Fatalf("TestGetConn expected equality,\nexpected writer[%v],\nconn writer[%v]", expect.Writer, conn.Writer)
	}
}

package model

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.chat/utils"
)

func TestIsConnEmpty(t *testing.T) {
	t.Logf("TestIsConnEmpty started")
	uc := make(UserConn, 0)
	isConn, conn := uc.IsConn("test_user")
	if isConn {
		t.Errorf("TestIsConnEmpty user was not supposed to be conn [%s]", conn.Log())
	} else if conn != nil {
		t.Errorf("TestIsConnEmpty expected empty, got [%s]", conn.Log())
	}
}

func TestIsConn(t *testing.T) {
	t.Logf("TestIsConn started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}
	uc.Add("test_user", "test_origin", w, *r)
	isConn, conn := uc.IsConn("test_user")
	if !isConn {
		t.Errorf("TestIsConn user was supposed to be conn")
	} else if conn == nil {
		t.Errorf("TestIsConn expected to have connection")
	}
}

func TestAdd(t *testing.T) {
	t.Logf("TestAdd started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}

	conn := uc.Add("test_user", "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestAdd expected conn, got NIL")
	} else if conn.User != "test_user" || conn.Origin != "test_origin" {
		t.Errorf("TestAdd unexpected user and origin [%s] [%s]", conn.User, conn.Origin)
	}
}

func TestGetEmpty(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := UserConn{}
	conn, err := uc.Get("test-req-id", "test_user")
	if conn != nil {
		t.Errorf("TestGetEmpty expected empty, got [%s]", conn.Log())
	} else if err == nil {
		t.Errorf("TestGetEmpty expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), "test_user") {
		t.Errorf("TestGetEmpty unexpected error [%s]", err)
	}
}

func TestGet(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := UserConn{}
	conn, err := uc.Get("test-req-id", "test_user")
	if conn != nil {
		t.Errorf("TestGetEmpty expected empty, got [%s]", conn.Log())
	} else if err == nil {
		t.Errorf("TestGetEmpty expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), "test_user") {
		t.Errorf("TestGetEmpty unexpected error [%s]", err)
	}
}

func TestDropEmpty(t *testing.T) {
	t.Logf("TestDropEmpty started")
	uc := UserConn{}
	conn := Conn{}
	uc.Drop("test-req-id", &conn)
}

func TestDrop(t *testing.T) {
	t.Logf("TestDrop started")
	r, err := http.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	utils.SetReqId(r, &reqId)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	uc := make(UserConn, 0)
	testUser := "test_user"
	addedConn := uc.Add(testUser, "test_origin", w, *r)
	if addedConn == nil {
		t.Errorf("failed to add conn")
	}

	conn, err := uc.Get(reqId, "test_user")
	if err != nil || conn == nil {
		t.Errorf("TestDrop unexpected error, %s", err)
	}

	uc.Drop(reqId, conn)
	conn, err = uc.Get(reqId, "test_user")
	if conn != nil {
		t.Errorf("TestDrop expected empty, got [%s]", conn.Log())
	} else if err == nil {
		t.Errorf("TestDrop expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), testUser) {
		t.Errorf("TestDrop unexpected error [%s]", err)
	}
}

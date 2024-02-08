package model

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetEmpty(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := UserConn{}
	conn, err := uc.Get("test_user")
	if conn != nil {
		t.Errorf("TestGetEmpty expected empty, got [%s]", conn.Log())
	} else if err == nil {
		t.Errorf("TestGetEmpty expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), "test_user") {
		t.Errorf("TestGetEmpty unexpected error [%s]", err)
	}
}

func TestAdd(t *testing.T) {
	t.Logf("TestAdd started")
	r, err := http.NewRequest("GET", "/some-route", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	uc := UserConn{}
	uc.Add("test_user", "test_origin", w, *r)
	conn, err := uc.Get("test_user")
	if err != nil {
		t.Errorf("TestAdd expected no error, got [%s]", err)
	} else if conn == nil {
		t.Errorf("TestAdd expected conn, got NIL")
	} else if conn.User != "test_user" || conn.Origin != "test_origin" {
		t.Errorf("TestAdd unexpected user and origin [%s] [%s]", conn.User, conn.Origin)
	}
}

func TestDropEmpty(t *testing.T) {
	t.Logf("TestDropEmpty started")
	uc := UserConn{}
	uc.Drop("test_user")
}

func TestDrop(t *testing.T) {
	t.Logf("TestDrop started")
	r, err := http.NewRequest("GET", "/some-route", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	uc := UserConn{}
	uc.Add("test_user", "test_origin", w, *r)
	conn, err := uc.Get("test_user")
	if err != nil || conn == nil {
		t.Errorf("TestDrop unexpected error, %s, expected NIL but got [%s]", err, conn.Log())
	}

	uc.Drop("test_user")
	conn, err = uc.Get("test_user")
	if conn != nil {
		t.Errorf("TestDrop expected empty, got [%s]", conn.Log())
	} else if err == nil {
		t.Errorf("TestDrop expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), "test_user") {
		t.Errorf("TestDrop unexpected error [%s]", err)
	}
}

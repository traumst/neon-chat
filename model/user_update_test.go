package model

import (
	"strings"
	"testing"
)

func TestDefaultUserUpdate(t *testing.T) {
	t.Logf("TestDefaultUserUpdate started")
	up := UserUpdate{}
	if up.Type != UnknownUpdate {
		t.Errorf("TestDefaultUserUpdate expected UnknownUpdate, got [%d]", up.Type)
	} else if up.Chat != nil || up.Msg != nil {
		t.Errorf("TestDefaultUserUpdate expected only null values, but chat[%t],msg[%t]", up.Chat == nil, up.Msg == nil)
	} else if up.User != "" {
		t.Errorf("TestDefaultUserUpdate expected empty user, got [%s]", up.User)
	}
}

func TestDefaultUserUpdateLog(t *testing.T) {
	t.Logf("TestDefaultUserUpdateLog started")
	up := UserUpdate{}
	log := up.Log()
	conditions := strings.Contains(log, "UserUpdate") &&
		strings.Contains(log, "User") &&
		strings.Contains(log, "Type") &&
		strings.Contains(log, "nil")
	if !conditions {
		t.Errorf("TestDefaultUserUpdateLog unexpected default log [%s]", log)
	}
}

func TestInitUserUpdate(t *testing.T) {
	t.Logf("TestInitUserUpdate started")
	up := UserUpdate{User: "test_user_fdyguf"}
	log := up.Log()
	conditions := strings.Contains(log, "UserUpdate") &&
		strings.Contains(log, "User") &&
		strings.Contains(log, "test_user_fdyguf")
	if !conditions {
		t.Errorf("TestInitUserUpdate unexpected log [%s]", log)
	}
}

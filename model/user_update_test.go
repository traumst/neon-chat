package model

import (
	"fmt"
	"strings"
	"testing"
)

func TestDefaultUserUpdate(t *testing.T) {
	t.Logf("TestDefaultUserUpdate started")
	up := UserUpdate{}
	if up.Type != UnknownUpdate {
		t.Errorf("TestDefaultUserUpdate expected UnknownUpdate, got [%s]", up.Type.String())
	} else if up.Msg != "" {
		t.Errorf("TestDefaultUserUpdate expected only \"\" values, got msg[%s]", up.Msg)
	} else if up.User != "" {
		t.Errorf("TestDefaultUserUpdate expected empty user, got [%s]", up.User)
	}
}

func TestInitUserUpdate(t *testing.T) {
	t.Logf("TestInitUserUpdate started")
	up := UserUpdate{User: "test_user_fdyguf"}
	log := fmt.Sprintf("%+v", up)
	conditions := strings.Contains(log, "UnknownUpdate") &&
		strings.Contains(log, "test_user_fdyguf")
	if !conditions {
		t.Errorf("TestInitUserUpdate unexpected log [%s]", log)
	}
}

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
	} else if up.RawHtml != "" {
		t.Errorf("TestDefaultUserUpdate expected only \"\" values, got msg[%s]", up.RawHtml)
	} else if up.Author != "" {
		t.Errorf("TestDefaultUserUpdate expected empty user, got [%s]", up.Author)
	}
}

func TestInitUserUpdate(t *testing.T) {
	t.Logf("TestInitUserUpdate started")
	up := UserUpdate{Author: "test_user_fdyguf"}
	log := fmt.Sprintf("%+v", up)
	conditions := strings.Contains(log, "unknown") &&
		strings.Contains(log, "test_user_fdyguf")
	if !conditions {
		t.Errorf("TestInitUserUpdate unexpected log [%s]", log)
	}
}

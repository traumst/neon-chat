package app

import "testing"

func TestDefaultChat(t *testing.T) {
	t.Logf("TestDefaultChat started")
	c := Chat{}
	if c.Id != 0 {
		t.Errorf("TestDefaultChat expected UnknownUpdate, got [%d]", c.Id)
	} else if c.Name != "" {
		t.Errorf("TestDefaultChat expected empty user, got [%s]", c.Name)
	} else if c.OwnerId != 0 {
		t.Errorf("TestDefaultChat expected empty owner, got [%d]", c.OwnerId)
	} else if c.OwnerName != "" {
		t.Errorf("TestDefaultChat expected empty owner, got [%s]", c.OwnerName)
	}
	t.Logf("TestDefaultChat finished")
}

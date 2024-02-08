package controller

import (
	"testing"

	"go.chat/model"
)

func TestDropEmpty(t *testing.T) {
	t.Logf("TestDropEmpty started")
	uc := model.UserConn{}
	uc.Drop("test_user")
}

package http

import (
	"net/http/httptest"
	"testing"
)

func TestReqId(t *testing.T) {
	r := httptest.NewRequest("GET", "/some-route", nil)
	testId := "test-req-id"
	setId := SetReqId(r, &testId)
	if setId != testId {
		t.Errorf("Expected testId[%s], setId[%s]", testId, setId)
	}
	gotId := GetReqId(r)
	if gotId != testId {
		t.Errorf("Expected testId[%s], gotId[%s]", testId, gotId)
	}
}

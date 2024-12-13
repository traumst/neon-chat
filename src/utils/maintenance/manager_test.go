package maintenance

import (
	"math/rand"
	"testing"
	"time"
)

const randMin = 2000
const randMax = maxConcurrentCount

func getRandCount() int {
	return rand.Intn(randMax-randMin) + randMin
}

func TestManagerDefaults(t *testing.T) {
	mm := &maintenanceManager{}
	if mm.activeUsers != 0 {
		t.Errorf("expected activeUsers 0, got [%d]", mm.activeUsers)
	}
	if mm.inMaintenance != 0 {
		t.Errorf("expected inMaintenance 0, got [%d]", mm.inMaintenance)
	}
	if mm.IsInMaintenance() {
		t.Errorf("expected IsInMaintenance false, got true")
	}
}

func TestRaiseFlag(t *testing.T) {
	mm := &maintenanceManager{}
	err := mm.RaiseFlag()
	if err != nil {
		t.Errorf("expected no error, got [%s]", err)
	}
	if !mm.IsInMaintenance() {
		t.Errorf("expected IsInMaintenance true, got false")
	}
	err = mm.RaiseFlag()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestClearFlag(t *testing.T) {
	mm := &maintenanceManager{}
	mm.RaiseFlag()
	mm.ClearFlag()
	if mm.IsInMaintenance() {
		t.Errorf("expected IsInMaintenance false, got true")
	}
}

func TestClearFlagConcurrent(t *testing.T) {
	mm := &maintenanceManager{}
	go mm.RaiseFlag()
	go mm.ClearFlag()
	if mm.IsInMaintenance() {
		t.Errorf("expected IsInMaintenance false, got true")
	}
}

func TestClearFlagRepeatCall(t *testing.T) {
	mm := &maintenanceManager{}
	go mm.RaiseFlag()
	go mm.ClearFlag()
	if mm.IsInMaintenance() {
		t.Errorf("expected IsInMaintenance false, got true")
	}
	go mm.ClearFlag()
	go mm.ClearFlag()
}

func TestWaitMaintenanceComplete(t *testing.T) {
	mm := &maintenanceManager{}
	mm.RaiseFlag()
	go func() {
		time.Sleep(200 * time.Millisecond)
		mm.ClearFlag()
	}()
	if !mm.WaitMaintenanceComplete(500 * time.Millisecond) {
		t.Errorf("expected true, got false")
	}
}

func TestWaitMaintenanceCompleteNotYet(t *testing.T) {
	mm := &maintenanceManager{}
	mm.RaiseFlag()
	if mm.WaitMaintenanceComplete(100 * time.Millisecond) {
		t.Errorf("expected false, got true")
	}
}

func TestIncrUserCount(t *testing.T) {
	mm := &maintenanceManager{}
	var err error
	var count int
	max := getRandCount()
	for i := 1; i < max; i++ {
		err = mm.IncrUserCount()
		if err != nil {
			t.Errorf("expected no error, got [%s]", err)
		}
		count = int(mm.activeUsers)
		if count != i {
			t.Errorf("expected activeUsers 1, got [%d]", mm.activeUsers)
		}
	}
}

func TestIncrUserCountOver(t *testing.T) {
	mm := &maintenanceManager{}
	var err error
	max := maxConcurrentCount
	for i := 0; i < max; i++ {
		_ = mm.IncrUserCount()
	}
	err = mm.IncrUserCount()
	if err == nil {
		t.Errorf("expected error due to too many users")
	}
	count := int(mm.activeUsers)
	if count != maxConcurrentCount {
		t.Errorf("expected [%d], got [%d]", maxConcurrentCount, count)
	}
}

func TestDecrUserCount(t *testing.T) {
	mm := &maintenanceManager{}
	max := getRandCount()
	for i := 0; i < max; i++ {
		_ = mm.IncrUserCount()
	}
	countIn := int(mm.activeUsers)
	if countIn != max {
		t.Errorf("expected [%d], got [%d]", max, countIn)
	}
	var countOut int
	for i := 1; i <= max; i++ {
		mm.DecrUserCount()
		countOut = int(mm.activeUsers)
		if countOut != max-i {
			t.Errorf("expected [%d], got [%d]", max-i, countOut)
		}
	}
}

func TestDecrUserCountOver(t *testing.T) {
	mm := &maintenanceManager{}
	max := getRandCount()
	for i := 0; i < max; i++ {
		_ = mm.IncrUserCount()
	}
	countIn := int(mm.activeUsers)
	if countIn != max {
		t.Errorf("expected [%d], got [%d]", max, countIn)
	}
	for i := 1; i <= max; i++ {
		mm.DecrUserCount()
	}

	defer func() { _ = recover() }()
	mm.DecrUserCount()
	t.Errorf("expected panic, got none")
}

func TestWaitUsersLeave(t *testing.T) {
	mm := &maintenanceManager{}
	count := mm.WaitUsersLeave(10 * time.Millisecond)
	if count != 0 {
		t.Errorf("expected 0, got [%d]", count)
	}
	_ = mm.IncrUserCount()
	go func() {
		time.Sleep(500 * time.Millisecond)
		mm.DecrUserCount()
	}()
	count = mm.WaitUsersLeave(1 * time.Second)
	if count != 0 {
		t.Errorf("expected 0, got [%d]", count)
	}
}

func TestWaitUsersLeaveNotLeft(t *testing.T) {
	mm := &maintenanceManager{}
	count := mm.WaitUsersLeave(10 * time.Millisecond)
	if count != 0 {
		t.Errorf("expected 0, got [%d]", count)
	}
	_ = mm.IncrUserCount()
	count = mm.WaitUsersLeave(500 * time.Millisecond)
	if count != 1 {
		t.Errorf("expected 1, got [%d]", count)
	}
}

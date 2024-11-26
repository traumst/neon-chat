package utils

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

var MaintenanceManager *maintenanceManager = &maintenanceManager{}

type maintenanceManager struct {
	activeUsers   int32
	inMaintenance int32
}

func (m *maintenanceManager) RaiseFlag() error {
	log.Println("TRACE RaiseFlag called")
	if atomic.LoadInt32(&m.inMaintenance) == 1 {
		return fmt.Errorf("maintenance flag is already up")
	}
	atomic.StoreInt32(&m.inMaintenance, 1)
	return nil
}

func (m *maintenanceManager) ClearFlag() {
	log.Println("TRACE ClearFlag called")
	atomic.StoreInt32(&m.inMaintenance, 0)
}

func (m *maintenanceManager) IsInMaintenance() bool {
	log.Println("TRACE IsInMaintenance called")
	return atomic.LoadInt32(&m.inMaintenance) == 1
}

func (m *maintenanceManager) IncrUserCount() error {
	log.Println("TRACE IncrUserCount called")
	if atomic.LoadInt32(&m.inMaintenance) == 1 {
		return fmt.Errorf("maintenance flag is up")
	}
	atomic.AddInt32(&m.activeUsers, 1)
	return nil
}

func (m *maintenanceManager) DecrUserCount() {
	log.Println("TRACE DecrUserCount called")
	atomic.AddInt32(&m.activeUsers, -1)
}

func (m *maintenanceManager) WaitUsersLeave(timeout time.Duration) bool {
	log.Println("TRACE WaitUsersLeave called, current count", atomic.LoadInt32(&m.activeUsers))
	probe := time.NewTicker(100 * time.Millisecond)
	abort := time.NewTicker(timeout)
	defer func() {
		probe.Stop()
		abort.Stop()
	}()

	for {
		select {
		case <-probe.C:
			if atomic.LoadInt32(&m.activeUsers) == 0 {
				return true
			}
		case <-abort.C:
			return false
		}
	}
}

package maintenance

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

const maxConcurrentCount = 3000

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
	//log.Println("TRACE IsInMaintenance called")
	return atomic.LoadInt32(&m.inMaintenance) == 1
}

func (m *maintenanceManager) WaitMaintenanceComplete(timeout time.Duration) bool {
	log.Println("TRACE WaitMaintenanceComplete called, current count", atomic.LoadInt32(&m.activeUsers))
	probe := time.NewTicker(100 * time.Millisecond)
	abort := time.NewTicker(timeout)
	defer func() {
		probe.Stop()
		abort.Stop()
	}()

	for {
		select {
		case <-probe.C:
			if atomic.LoadInt32(&m.inMaintenance) == 0 {
				return true
			}
		case <-abort.C:
			return false
		}
	}
}

func (m *maintenanceManager) IncrUserCount() error {
	log.Println("TRACE IncrUserCount called")
	if atomic.LoadInt32(&m.inMaintenance) == 1 {
		return fmt.Errorf("maintenance flag is up")
	}
	count := atomic.LoadInt32(&m.activeUsers)
	if count >= maxConcurrentCount {
		return fmt.Errorf("max concurrent users reached")
	}
	atomic.AddInt32(&m.activeUsers, 1)
	return nil
}

func (m *maintenanceManager) DecrUserCount() {
	log.Println("TRACE DecrUserCount called")
	if atomic.LoadInt32(&m.activeUsers) == 0 {
		panic("count was already 0 on call to decrement")
	}
	atomic.AddInt32(&m.activeUsers, -1)
}

func (m *maintenanceManager) WaitUsersLeave(timeout time.Duration) int {
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
				return 0
			}
		case <-abort.C:
			return int(atomic.LoadInt32(&m.activeUsers))
		}
	}
}

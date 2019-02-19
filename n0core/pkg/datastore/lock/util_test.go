package lock

import (
	"testing"
	"time"
)

func TestWaitUntilLock(t *testing.T) {
	m := NewMemoryMutexTable(100)

	if !WaitUntilLock(m, "test", 300*time.Millisecond) {
		t.Errorf("failed to lock")
	}
	if WaitUntilLock(m, "test", 300*time.Millisecond) {
		t.Errorf("locked after locked")
	}
}

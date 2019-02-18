package lock

import "testing"

func TestMemoryMutexTable(t *testing.T) {
	mt := NewMemoryMutexTable(10000)

	if mt.IsLocked("test") {
		t.Errorf("precondition was locked")
	}
	if mt.Unlock("test") {
		t.Errorf("failed to unlock on precondition")
	}

	if !mt.Lock("test") {
		t.Errorf("failed to lock")
	}
	if !mt.IsLocked("test") {
		t.Errorf("is not locked after locked")
	}
	if mt.Lock("test") {
		t.Errorf("lock after locked")
	}

	if !mt.Unlock("test") {
		t.Errorf("failed to unlock")
	}
	if mt.IsLocked("test") {
		t.Errorf("is not unlocked after unlocked")
	}
}

func BenchmarkLock(b *testing.B) {
	mt := NewMemoryMutexTable(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go mt.Lock("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.Unlock("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
		go mt.IsLocked("test")
	}
}

package lock

type MutexTable interface {
	// Lock key
	// return that lock is succeeded
	Lock(key string) bool

	// Unlock key
	// return that unlock is succeeded
	Unlock(key string) bool

	// IsLocked return that key is locked
	IsLocked(key string) bool
}

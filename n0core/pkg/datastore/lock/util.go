package lock

import "time"

func WaitUntilLock(mutex MutexTable, key string, timeout time.Duration) bool {
	over := time.After(timeout)

	for {
		select {
		case <-over:
			return false

		default:
			if mutex.Lock(key) {
				return true
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

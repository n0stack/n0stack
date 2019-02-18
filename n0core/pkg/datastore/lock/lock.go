package lock

type MemoryMutexTable struct {
	table   map[string]bool
	request chan mutexRequest
}

func NewMemoryMutexTable(requestBuffer int) *MemoryMutexTable {
	mt := &MemoryMutexTable{
		table:   make(map[string]bool),
		request: make(chan mutexRequest, requestBuffer),
	}

	go mt.mutexThread()

	return mt
}

func (mt *MemoryMutexTable) Lock(key string) bool {
	ch := make(chan mutexResult)
	defer close(ch)

	mt.request <- mutexRequest{
		Key:    key,
		Action: lock,
		Result: ch,
	}

	for r := range ch {
		return r.Succeeded
	}

	return false
}

func (mt *MemoryMutexTable) Unlock(key string) bool {
	ch := make(chan mutexResult)
	defer close(ch)

	mt.request <- mutexRequest{
		Key:    key,
		Action: unlock,
		Result: ch,
	}

	for r := range ch {
		return r.Succeeded
	}

	return false
}

func (mt MemoryMutexTable) IsLocked(key string) bool {
	ch := make(chan mutexResult)
	defer close(ch)

	mt.request <- mutexRequest{
		Key:    key,
		Action: isLocked,
		Result: ch,
	}

	for r := range ch {
		return r.Locked
	}

	return false
}

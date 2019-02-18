package lock

type MutexTable struct {
	table   map[string]bool
	request chan mutexRequest
}

func NewMutexTable(requestBuffer int) *MutexTable {
	mt := &MutexTable{
		table:   make(map[string]bool),
		request: make(chan mutexRequest, requestBuffer),
	}

	go mt.mutexThread()

	return mt
}

func (mt *MutexTable) Lock(key string) bool {
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

func (mt *MutexTable) Unlock(key string) bool {
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

func (mt MutexTable) IsLocked(key string) bool {
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

package raceutil

import (
	"io"
	"sync"
)

// https://www.reddit.com/r/golang/comments/6bpmtj/writing_to_a_file_on_multiple_threads/
type LockedWriter struct {
	m *sync.Mutex
	w io.Writer
}

func NewLockedWriter(w io.Writer) *LockedWriter {
	return &LockedWriter{
		m: &sync.Mutex{},
		w: w,
	}
}

func (lw *LockedWriter) Write(b []byte) (n int, err error) {
	lw.m.Lock()
	defer lw.m.Unlock()
	return lw.w.Write(b)
}

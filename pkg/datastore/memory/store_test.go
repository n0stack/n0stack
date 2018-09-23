package memory

import (
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
)

func TestMemoryDatastore(t *testing.T) {
	m := NewMemoryDatastore()

	if err := m.Apply("test", &empty.Empty{}); err != nil {
		t.Errorf("Failed to apply empty: err='%s'", err.Error())
	}

	e := &empty.Empty{}
	if err := m.Get("test", e); err != nil {
		t.Errorf("Failed to get empty: err='%s'", err.Error())
	} else if e == nil {
		t.Errorf("Failed to get empty: result is nil")
	}

	if err := m.Delete("test"); err != nil {
		t.Errorf("Failed to delete: err='%s'", err.Error())
	}
}

// +build !small

package iproute2

import "testing"

func TestBridge(t *testing.T) {
	b, err := NewBridge("test")
	if err != nil {
		t.Fatalf("Failed to create bridge: err='%s'", err.Error())
	}

	if err := b.Up(); err != nil {
		t.Errorf("Failed to up bridge: err='%s'", err.Error())
	}

	if err := b.Delete(); err != nil {
		t.Errorf("Failed to delete bridge: err='%s'", err.Error())
	}
}

func TestExistingBridge(t *testing.T) {
	b, err := NewBridge("test")
	if err != nil {
		t.Fatalf("Failed to create bridge: err='%s'", err.Error())
	}
	defer b.Delete()

	if _, err := NewBridge("test"); err != nil {
		t.Fatalf("Failed to find existing bridge: err='%s'", err.Error())
	}
}

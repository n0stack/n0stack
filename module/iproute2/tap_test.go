package iproute2

import "testing"

func TestTap(t *testing.T) {
	nt, err := NewTap("test")
	if err != nil {
		t.Fatalf("Failed to create tap: err='%s'", err.Error())
	}

	if err := nt.Up(); err != nil {
		t.Errorf("Failed to up tap: err='%s'", err.Error())
	}

	if err := nt.Delete(); err != nil {
		t.Errorf("Failed to delete tap: err='%s'", err.Error())
	}
}

func TestExistingTap(t *testing.T) {
	nt, err := NewTap("test")
	if err != nil {
		t.Fatalf("Failed to create tap: err='%s'", err.Error())
	}
	defer nt.Delete()

	if _, err := NewTap("test"); err != nil {
		t.Errorf("Failed to find existing tap: err='%s'", err.Error())
	}
}

func TestTapSetMasterAsBridge(t *testing.T) {
	b, err := NewBridge("br-test")
	if err != nil {
		t.Fatalf("Failed to create bridge: err='%s'", err.Error())
	}
	defer b.Delete()

	nt, err := NewTap("tap-test")
	if err != nil {
		t.Fatalf("Failed to create tap: err='%s'", err.Error())
	}
	defer nt.Delete()

	if err := nt.SetMaster(b); err != nil {
		t.Errorf("Failed to set master: err='%s'", err.Error())
	}
}

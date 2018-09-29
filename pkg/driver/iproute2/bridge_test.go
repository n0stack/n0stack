// +build medium
// +build !without_root

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

	if err := b.SetAddress("10.255.255.1/24"); err != nil {
		t.Errorf("Failed to set address: err='%s'", err.Error())
	}

	if _, err := b.GetIPv4(); err != nil {
		t.Errorf("Failed to get address: err='%s'", err.Error())
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

func TestListSlaves(t *testing.T) {
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

	if l, err := b.ListSlaves(); err != nil {
		t.Errorf("Failed to list slaves: err='%s'", err.Error())
	} else if len(l) != 1 {
		t.Errorf("Wrong the number of slaves: want='%d', have='%d'", 1, len(l))
	} else if l[0] != nt.Name() {
		t.Errorf("Got wrong link about index: want='%s', have='%s'", nt.Name(), l[0])
	}
}

// +build medium
// +build !without_root

package qemu

import (
	"net"
	"os"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"n0st.ac/n0stack/n0core/pkg/driver/iproute2"
	uuid "github.com/satori/go.uuid"
)

func TestAttachTap(t *testing.T) {
	b, err := iproute2.NewBridge("br-test")
	if err != nil {
		t.Fatalf("Failed to create bridge: err='%s'", err.Error())
	}
	defer b.Delete()

	i, err := iproute2.NewTap("tap-test")
	if err != nil {
		t.Fatalf("Failed to create tap: err='%s'", err.Error())
	}
	defer i.Delete()

	i.SetMaster(b)

	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu("test")
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Delete()

	if _, ok := os.LookupEnv("DISABLE_KVM"); ok {
		q.isKVM = false
	}

	m, _ := bytefmt.ToBytes("512M")
	if err := q.Start(id, "monitor.sock", 10000, 1, m); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	hw, _ := net.ParseMAC("52:54:01:23:45:67")
	if err := q.AttachTap("test", i.Name(), hw); err != nil {
		t.Errorf("Failed to attach tap: err='%s'", err.Error())
	}

	if err := q.Boot(); err != nil {
		t.Errorf("Failed to boot: err='%s'", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusRunning {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusRunning, s)
	}
}

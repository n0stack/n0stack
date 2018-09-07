package qemu

import (
	"testing"

	"code.cloudfoundry.org/bytefmt"

	uuid "github.com/satori/go.uuid"
)

func TestStartProcess(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	if q.IsRunning() {
		t.Fatalf("Test environment is invalid, process is already existing: uuid='%s'", id.String())
	}

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Errorf("Failed to start process: err='%s'", err.Error())
	}
	if !q.IsRunning() {
		t.Errorf("Failed to start process, qemu is not running yet")
	}

	if err := q.Kill(); err != nil {
		t.Errorf("Failed to kill process: err='%s'", err.Error())
	}
	if q.IsRunning() {
		t.Errorf("Failed to kill process, qemu is running yet")
	}
}

func TestBoot(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Kill()

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusPreLaunch {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusPreLaunch, s)
	}

	if err := q.Boot(); err != nil {
		t.Errorf("Failed to boot: err='%s'", err.Error())
	}

	s, err = q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusRunning {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusRunning, s)
	}
}

func TestReset(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Kill()

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusPreLaunch {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusPreLaunch, s)
	}

	if err := q.Boot(); err != nil {
		t.Errorf("Failed to boot: err='%s'", err.Error())
	}
	if err := q.Reset(); err != nil {
		t.Errorf("Failed to reset: err='%s'", err.Error())
	}

	s, err = q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusRunning {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusRunning, s)
	}
}

// func TestShutdown(t *testing.T) {
// 	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
// 	q, err := OpenQemu(&id)
// 	if err != nil {
// 		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
// 	}
// 	if q.IsRunning() {
// 		t.Fatalf("Test environment is invalid, process is already existing: uuid='%s'", id.String())
// 	}

// 	b, _ := bytefmt.ToBytes("512M")
// 	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
// 		t.Errorf("Failed to start process: err='%s'", err.Error())
// 	}
// 	if !q.IsRunning() {
// 		t.Errorf("Failed to start process, qemu is not running yet")
// 	}

// 	s, err := q.Status()
// 	if err != nil {
// 		t.Errorf("Failed to get status: err='%s'", err.Error())
// 	}
// 	if s != StatusPreLaunch {
// 		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusPreLaunch, s)
// 	}

// 	if err := q.Boot(); err != nil {
// 		t.Errorf("Failed to boot: err='%s'", err.Error())
// 	}
// 	if err := q.Shutdown(); err != nil {
// 		t.Errorf("Failed to reset: err='%s'", err.Error())
// 	}

// 	s, err = q.Status()
// 	if err != nil {
// 		t.Errorf("Failed to get status: err='%s'", err.Error())
// 	}
// 	if s != StatusRunning {
// 		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusRunning, s)
// 	}

// 	if err := q.Kill(); err != nil {
// 		t.Errorf("Failed to kill process: err='%s'", err.Error())
// 	}
// 	if q.IsRunning() {
// 		t.Errorf("Failed to kill process, qemu is running yet")
// 	}
// }

package qemu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/digitalocean/go-qemu/qmp"
	"github.com/digitalocean/go-qemu/qmp/raw"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
)

type Qemu struct {
	proc *process.Process

	// args
	name    string
	qmpPath string
	isKVM   bool

	qmp qmp.Monitor
	m   *raw.Monitor
}

func OpenQemu(name string) (*Qemu, error) {
	q := &Qemu{
		name:  name,
		isKVM: true,
	}

	if err := q.init(); err != nil {
		return nil, err
	}

	return q, nil
}

func (q *Qemu) Close() error {
	if q.m == nil {
		return nil
	}

	return q.qmp.Disconnect()
}

// func (q Qemu) Reset() error {
// 	return nil
// }

func (q Qemu) HardReset() error {
	return q.m.SystemReset()
}

func (q Qemu) Shutdown() error {
	return q.m.SystemPowerdown()
}

func (q Qemu) HardShutdown() error {
	if err := q.m.Stop(); err != nil {
		return err
	}

	return q.m.SystemReset()
}

func (q Qemu) Boot() error {
	if s, err := q.Status(); err != nil {
		return err
	} else if s == StatusShutdown {
		if err := q.m.SystemReset(); err != nil {
			return err
		}
	}

	if err := q.m.Cont(); err != nil {
		return err
	}

	return q.m.SystemWakeup()
}

type Status int

const (
	StatusDebug         Status = Status(raw.RunStateDebug)
	StatusFinishMigrate Status = Status(raw.RunStateFinishMigrate)
	StatusGuestPanicked Status = Status(raw.RunStateGuestPanicked)
	StatusIOError       Status = Status(raw.RunStateIOError)
	StatusInMigrate     Status = Status(raw.RunStateInmigrate)
	StatusInternalError Status = Status(raw.RunStateInternalError)
	StatusPaused        Status = Status(raw.RunStatePaused)
	StatusPostMigrate   Status = Status(raw.RunStatePostmigrate)
	StatusPreLaunch     Status = Status(raw.RunStatePrelaunch)
	StatusRestoreVM     Status = Status(raw.RunStateRestoreVM)
	StatusRunning       Status = Status(raw.RunStateRunning)
	StatusSaveVM        Status = Status(raw.RunStateSaveVM)
	StatusShutdown      Status = Status(raw.RunStateShutdown)
	StatusSuspended     Status = Status(raw.RunStateSuspended)
	StatusWatchdog      Status = Status(raw.RunStateWatchdog)
)

func (q Qemu) Status() (Status, error) {
	s, err := q.m.QueryStatus()
	if err != nil {
		if strings.Contains(err.Error(), "not running") {
			return StatusShutdown, nil
		}

		return 0, err
	}

	return Status(s.Status), nil
}

func (q Qemu) IsRunning() bool {
	if q.proc == nil {
		return false
	}

	return true
}

func (q *Qemu) Delete() error {
	if q.proc != nil {
		if err := q.proc.Kill(); err != nil {
			return fmt.Errorf("Failed to kill process: err='%s'", err.Error())
		}

		q.proc = nil
	}

	if err := os.Remove(q.qmpPath); err != nil {
		return fmt.Errorf("Failed to delete QMP socket: err='%s'", err.Error())
	}

	return nil
}

func (q *Qemu) Start(id uuid.UUID, qmpPath string, vncWebsocketPort uint16, vcpus uint32, memory uint64) error {
	qmpPath, err := filepath.Abs(qmpPath)
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of qmpPath: err='%s'", err.Error())
	}

	args := []string{
		"qemu-system-x86_64",

		// -- QEMU metadata --
		"-uuid",
		id.String(),
		"-name",
		fmt.Sprintf("guest=%s,debug-threads=on", q.name),
		"-msg",
		"timestamp=on",

		// Config
		"-daemonize",
		"-nodefaults",     // Don't create default devices
		"-no-user-config", // The "-no-user-config" option makes QEMU not load any of the user-provided config files on sysconfdir
		"-S",              // Do not start CPU at startup
		"-no-shutdown",    // Don't exit QEMU on guest shutdown
		"-global",
		"kvm-pit.lost_tick_policy=discard",
		// "-pidfile",
		// "",

		// QMP
		"-chardev",
		fmt.Sprintf("socket,id=charmonitor,path=%s,server,nowait", qmpPath),
		"-mon",
		"chardev=charmonitor,id=monitor,mode=control",

		// -- BIOS --
		// boot priority
		"-boot",
		"menu=on,strict=on",

		// keyboard
		"-k",
		"en-us",

		// VNC
		"-vnc",
		fmt.Sprintf("0.0.0.0:%d,websocket=%d", GetNewListenPort(5900)-5900, vncWebsocketPort),

		// clock
		"-rtc",
		"base=utc,driftfix=slew",
		"-global",
		"kvm-pit.lost_tick_policy=delay",
		"-no-hpet",

		// CPU
		// TODO: 必要があればmonitorを操作してhotaddできるようにする
		// TODO: スケジューリングが可能かどうか調べる
		"-smp",
		fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", vcpus, vcpus),
		"-cpu",
		"host",
		"-enable-kvm",

		// Memory
		"-m",
		fmt.Sprintf("%s", bytefmt.ByteSize(memory)),
		// "-device",
		// "virtio-balloon-pci,id=balloon0,bus=pci.0", // dynamic configurations
		"-realtime",
		"mlock=off",

		// VGA controller
		"-device",
		"VGA,id=video0,bus=pci.0",

		// SCSI controller
		"-device",
		"lsi53c895a,bus=pci.0,id=scsi0",

		// Serial device
		"-device",
		"virtio-serial-pci,id=virtio-serial0,bus=pci.0",
		"-chardev",
		"pty,id=charserial0",
		"-device",
		"isa-serial,chardev=charserial0,id=serial0",
	}

	if !q.isKVM {
		// remove "-cpu", "host" and "-enable-kvm", because kvm is disable
		for i, v := range args {
			if v == "-cpu" {
				args = append(args[:i], args[i+3:]...)
			}
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil { // TODO: combine でもいいかもしれない
		return fmt.Errorf("Failed to start process: args='%s', out='%s', err='%s'", args, string(out), err.Error())
	}

	if err := q.init(); err != nil {
		return fmt.Errorf("Failed to initialize: args='%s', err='%s'", args, err.Error())
	}

	return nil
}

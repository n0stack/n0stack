package kvm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/digitalocean/go-qemu/qmp"
	"github.com/shirou/gopsutil/process"

	"code.cloudfoundry.org/bytefmt"
	"github.com/satori/go.uuid"

	"github.com/n0stack/n0core/lib"
	n0stack "github.com/n0stack/proto"
	"github.com/n0stack/proto/device/vm"
	"github.com/n0stack/proto/resource/cpu"
)

type (
	Agent struct {
		// DB *gorm.DB
	}

	kvm struct {
		id      uuid.UUID
		workDir string

		args   []string
		pid    int
		qmp    *qmp.SocketMonitor
		status *vm.Status
	}
)

const (
	modelType = "device/vm/kvm"
)

func (k kvm) getInstanceName(n string) string {
	return fmt.Sprintf("n0core-%s", n)
}

func getVM(model *n0stack.Model) (*kvm, *n0stack.Notification) {
	k := &kvm{}

	var err error
	k.id, err = uuid.FromBytes(model.Id)
	if err != nil {
		return nil, lib.MakeNotification("getVM.validateUUID", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	k.workDir, err = lib.GetWorkDir(modelType, k.id)
	if err != nil {
		return nil, lib.MakeNotification("getVM.getWorkDir", false, fmt.Sprintf("error message '%s', when creating work directory, '%s'", k.workDir, err.Error()))
	}

	ps, err := process.Processes()
	if err != nil {
		return nil, lib.MakeNotification("getVM.getProcessList", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	for _, p := range ps {
		c, _ := p.Cmdline() // エラーが発生する場合が考えられない
		// println(c)
		if strings.Contains(c, k.id.String()) {
			k.args, _ = p.CmdlineSlice()

			k.pid = int(p.Pid)
			return k, lib.MakeNotification("getVM", true, fmt.Sprintf("Already running: pid=%d", k.pid))
		}
	}

	return k, lib.MakeNotification("getVM", true, "Not running QEMU process")
}

func (k *kvm) runVM(spec *vm.Spec) *n0stack.Notification {
	switch spec.Cpu.Architecture {
	case cpu.Architecture_x86_64:
		k.args = []string{"qemu-system-x86_64"}
	}

	// -- QEMU metadata --
	k.args = append(k.args, "-uuid")
	k.args = append(k.args, k.id.String())
	k.args = append(k.args, "-name")
	k.args = append(k.args, fmt.Sprintf("guest=%s,debug-threads=on", k.getInstanceName(spec.Device.Model.Name)))
	k.args = append(k.args, "-msg")
	k.args = append(k.args, "timestamp=on")

	k.args = append(k.args, "-nodefaults")     // Don't create default devices
	k.args = append(k.args, "-no-user-config") // The "-no-user-config" option makes QEMU not load any of the user-provided config files on sysconfdir
	k.args = append(k.args, "-S")              // Do not start CPU at startup
	k.args = append(k.args, "-no-shutdown")    // Don't exit QEMU on guest shutdown

	// QMP
	const monitorFile = "monitor.sock"
	qmpPath := filepath.Join(k.workDir, monitorFile)
	k.args = append(k.args, "-chardev")
	k.args = append(k.args, fmt.Sprintf("socket,id=charmonitor,path=%s,server,nowait", qmpPath))
	k.args = append(k.args, "-mon")
	k.args = append(k.args, "chardev=charmonitor,id=monitor,mode=control")

	// -- BIOS --
	// boot priority
	k.args = append(k.args, "-boot")
	k.args = append(k.args, "menu=on,strict=on")

	// keyboard
	k.args = append(k.args, "-k")
	k.args = append(k.args, "en-us") // vm.Spec.Keymapみたいなので取得できるようにする

	// VNC
	k.args = append(k.args, "-vnc")
	k.args = append(k.args, ":0,websocket=5700") // ぶつからないようにポートを設定する必要がある, unix socketでも可 unix:$workdir/vnc.sock,websocket

	// clock
	k.args = append(k.args, "-rtc")
	k.args = append(k.args, "base=utc,driftfix=slew")
	k.args = append(k.args, "-global")
	k.args = append(k.args, "kvm-pit.lost_tick_policy=delay")
	k.args = append(k.args, "-no-hpet")

	// CPU
	// TODO: 必要があればmonitorを操作してhotaddできるようにする
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-cpu")
	k.args = append(k.args, "host")
	k.args = append(k.args, "-smp")
	k.args = append(k.args, fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", spec.Cpu.Vcpus, spec.Cpu.Vcpus))
	k.args = append(k.args, "-enable-kvm")
	// return true, "Succeeded to check cpu configurations"

	// Memory
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-m")
	k.args = append(k.args, fmt.Sprintf("%s", bytefmt.ByteSize(spec.Memory.Bytes)))
	k.args = append(k.args, "-device")
	k.args = append(k.args, "virtio-balloon-pci,id=balloon0,bus=pci.0") // dynamic configurations
	k.args = append(k.args, "-realtime")
	k.args = append(k.args, "mlock=off")

	// VGA controller
	k.args = append(k.args, "-device")
	k.args = append(k.args, "VGA,id=video0,bus=pci.0")

	// SCSI controller
	k.args = append(k.args, "-device")
	k.args = append(k.args, "lsi53c895a,bus=pci.0,id=scsi0")

	cmd := exec.Command(k.args[0], k.args[1:]...)
	if err := cmd.Start(); err != nil {
		return lib.MakeNotification("startQEMUProcess.startProcess", false, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args))
	}
	k.pid = cmd.Process.Pid

	// バックグラウンドプロセスにならなくなってしまう
	// ただ正常にプロセスが起動したかは待たないとわからない
	// done := make(chan error)
	// go func() {
	// 	done <- cmd.Wait()
	// }()

	// select {
	// case <-time.After(3 * time.Second):
	// 	return lib.MakeNotification("startQEMUProcess", true, "")
	// case err := <-done:
	// 	return lib.MakeNotification("startQEMUProcess.waitError", true, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args)) // stderrを表示できるようにする必要がある
	// }

	return lib.MakeNotification("startQEMUProcess", true, "")
}

func (k kvm) getQMPPath() string {
	chardev := map[string]string{}
	chardevID := ""

	for i, a := range k.args {
		switch {
		case a == "-mon":
			ops := strings.Split(k.args[i+1], ",")
			for _, o := range ops {
				if strings.HasPrefix(o, "chardev=") {
					chardevID = strings.Split(o, "=")[1]
				}
			}

		case a == "-chardev":
			var (
				id string
				p  string
			)

			ops := strings.Split(k.args[i+1], ",")
			for _, o := range ops {
				switch {
				case strings.HasPrefix(o, "id="):
					id = strings.Split(o, "=")[1]
				case strings.HasPrefix(o, "path="):
					p = strings.Split(o, "=")[1]
				}
			}

			chardev[id] = p
		}
	}

	return chardev[chardevID]
}

func (k *kvm) connectQMP() *n0stack.Notification {
	qmpPath := k.getQMPPath()

	var err error
	k.qmp, err = qmp.NewSocketMonitor("unix", qmpPath, 5*time.Second)
	if err != nil {
		return lib.MakeNotification("connectQMP", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return lib.MakeNotification("connectQMP", true, "")
}

func (k *kvm) bootVM() *n0stack.Notification {
	k.qmp.Connect()
	defer k.qmp.Disconnect()

	cmd := []byte(`{ "execute": "cont" }`)
	raw, err := k.qmp.Run(cmd)
	if err != nil {
		k.status.RunLevel = vm.RunLevel_SHUTDOWN
		return lib.MakeNotification("bootVM", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	k.status.RunLevel = vm.RunLevel_RUNNING
	return lib.MakeNotification("bootVM", true, fmt.Sprintf("qmp response '%s'", raw))
}

// func (t tap) getMACAddr() *net.HardwareAddr {
// 	c := crc32.ChecksumIEEE(t.id.Bytes())

// 	return &net.HardwareAddr{0x52, 0x54, c[2:]}
// }

// Apply スペックを元にステートレスに適用する
func (a *Agent) Apply(ctx context.Context, spec *vm.Spec) (*n0stack.Notification, error) {
	// ps auxfww | grep $uuid
	k, n := getVM(spec.Device.Model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		// check CPU usage
		// check Memory usage

		// qemu-system...
		n = k.runVM(spec)
		if !n.Success {
			return n, nil
		}
	}

	// qmp-shell .../monitor.sock
	n = k.connectQMP()
	if !n.Success {
		return n, nil
	}

	// (QEMU) cont
	// Applyした時に毎回ブートしなければいけないわけではない
	n = k.bootVM()
	return n, nil

	// (QEMU) ...
	// conn :=
	// vcl := volume.NewRepositoryClient(conn)
	// vcl.
	// k.attachVolume()
	// k.attachNIC()

	return lib.MakeNotification("Apply", true, ""), nil
}

func (k kvm) kill() *n0stack.Notification {
	p, _ := os.FindProcess(k.pid)
	if err := p.Kill(); err != nil {
		return lib.MakeNotification("Kill", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return lib.MakeNotification("Kill", true, "")
}

func (a *Agent) Delete(ctx context.Context, model *n0stack.Model) (*n0stack.Notification, error) {
	// ps auxfww | grep $uuid
	k, n := getVM(model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		return lib.MakeNotification("Delete", true, "Process is not existing"), nil
	}

	// kill $qemu
	n = k.kill()
	if !n.Success {
		return n, nil
	}

	return lib.MakeNotification("Delete", true, ""), nil
}

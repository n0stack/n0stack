package kvm

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/digitalocean/go-qemu/qmp"
	pkvm "github.com/n0stack/go.proto/kvm/v0"
	pnotification "github.com/n0stack/go.proto/notification/v0"
	"github.com/n0stack/n0core/lib"
	"github.com/n0stack/n0core/notification"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
)

type kvm struct {
	pkvm.Spec
	pkvm.Status

	id      uuid.UUID
	workDir string

	args []string
	pid  int
	qmp  *qmp.SocketMonitor
}

// ps auxfww | grep $uuid
func (k *kvm) getKVM() *pnotification.Notification {
	var err error
	k.workDir, err = lib.GetWorkDir(modelType, k.id)
	if err != nil {
		return notification.MakeNotification("getVM.GetWorkDir", false, fmt.Sprintf("error message '%s', when creating work directory, '%s'", k.workDir, err.Error()))
	}

	ps, err := process.Processes()
	if err != nil {
		return notification.MakeNotification("getVM.getProcessList", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	for _, p := range ps {
		c, _ := p.Cmdline()                                                     // エラーが発生する場合が考えられない
		if strings.Contains(c, k.id.String()) && strings.HasPrefix(c, "qemu") { // このfilterはガバガバなので後でリファクタリングする
			k.args, _ = p.CmdlineSlice()

			k.pid = int(p.Pid)
			return notification.MakeNotification("getVM", true, fmt.Sprintf("Already running: pid=%d", k.pid))
		}
	}

	return notification.MakeNotification("getVM", true, "Not running QEMU process")
}

// qemu-system...
func (k *kvm) runVM(vcpus uint32, memory uint64) *pnotification.Notification {
	k.args = []string{"qemu-system-x86_64"}

	// -- QEMU metadata --
	k.args = append(k.args, "-uuid")
	k.args = append(k.args, k.id.String())
	k.args = append(k.args, "-name")
	k.args = append(k.args, fmt.Sprintf("guest=%s,debug-threads=on", k.getInstanceName()))
	k.args = append(k.args, "-msg")
	k.args = append(k.args, "timestamp=on")

	k.args = append(k.args, "-daemonize")
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
	k.args = append(k.args, fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", vcpus, vcpus))
	k.args = append(k.args, "-enable-kvm")
	// return true, "Succeeded to check cpu configurations"

	// Memory
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-m")
	k.args = append(k.args, fmt.Sprintf("%s", bytefmt.ByteSize(memory)))
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
		return notification.MakeNotification("startQEMUProcess.startProcess", false, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args))
	}
	k.pid = cmd.Process.Pid

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(1 * time.Second):
		return notification.MakeNotification("startQEMUProcess", true, "")
	case err := <-done:
		if err != nil {
			return notification.MakeNotification("startQEMUProcess.waitError", false, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args)) // stderrを表示できるようにする必要がある
		}
		return notification.MakeNotification("startQEMUProcess.wait", true, fmt.Sprintf("args '%s'", k.args)) // stderrを表示できるようにする必要がある

	}
}

// qmp-shell .../monitor.sock
// TODO: 他のプロセスがソケットにつなげていた場合、何故か無制限にロックされてしまう
func (k *kvm) connectQMP() *pnotification.Notification {
	qmpPath := k.getQMPPath()

	var err error
	k.qmp, err = qmp.NewSocketMonitor("unix", qmpPath, 2*time.Second)
	if err != nil {
		return notification.MakeNotification("connectQMP", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return notification.MakeNotification("connectQMP", true, "")
}

// (QEMU) cont
func (k *kvm) bootVM() *pnotification.Notification {
	cmd := []byte(`{ "execute": "cont" }`)
	raw, err := k.qmp.Run(cmd)
	if err != nil {
		k.RunLevel = pkvm.Status_SHUTDOWN
		return notification.MakeNotification("bootVM", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	k.RunLevel = pkvm.Status_RUNNING
	return notification.MakeNotification("bootVM", true, fmt.Sprintf("qmp response '%s'", raw))
}

// // // (QEMU) query-block
// // func (k kvm) listVolumes() *n0stack.Notification {

// // 	return lib.MakeNotification("listVolumes", true, "")
// // }

// (QEMU) blockdev-add options={"driver":"qcow2","id":"drive-virtio-disk0","file":{"driver":"file","filename":"/home/h-otter/wk/test-qemu/ubuntu16.04.qcow2"}}
// (QEMU) device_add driver=virtio-blk-pci bus=pci.0 scsi=off drive=drive-virtio-disk0 id=virtio-disk0 bootindex=1
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - bootindexがどうやって更新されるのかがわからない
func (k *kvm) attachVolume(id, volumeUrl string, index int) *pnotification.Notification {
	u, err := url.Parse(volumeUrl)
	if err != nil {
		return notification.MakeNotification("attachVolume.parseURL", false, fmt.Sprintf("error message '%s', URL '%s', id '%s'", err.Error(), volumeUrl, id))
	}

	var cmd []byte
	switch {
	case u.Scheme == "file":
		cmd = []byte(fmt.Sprintf(`
				{
					"execute": "blockdev-add",
					"arguments": {
						"options": {
							"driver": "qcow2",
							"id": "drive-%s",
							"file": {
								"driver": "file",
								"filename": "%s"
							}
						}
					}
				}
			`, id, u.Path))
	}

	raw, err := k.qmp.Run(cmd)
	if err != nil && false { // already existsが発行されてしまう
		return notification.MakeNotification("attachVolume.blockdev-add", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	cmd = []byte(fmt.Sprintf(`
				{
					"execute": "device_add",
					"arguments": {
						"driver": "virtio-blk-pci",
						"id": "virtio-blk-%s",
						"drive": "drive-%s",
						"bus": "pci.0",
						"scsi": "off",
						"bootindex": "%d"
					}
				}
			`, id, id, index)) // bootindexはcdのために1を追加する

	raw, err = k.qmp.Run(cmd)
	if err != nil && false { // already existsが発行されてしまう
		return notification.MakeNotification("attachVolume.device_add", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	return notification.MakeNotification("attachVolume", true, "")
}

// (QEMU) netdev_add id=tap0 type=tap vhost=true ifname=tap0 script=no downscript=no
// (QEMU) device_add driver=virtio-net-pci netdev=tap0 id=test0 mac=52:54:00:df:89:29 bus=pci.0
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - MACアドレスを変更する
func (k *kvm) attachNIC(id, tap, mac string) *pnotification.Notification {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return notification.MakeNotification("attachNICs.ParseMAC", false, fmt.Sprintf("error message '%s', HWAddr '%s', id '%s'", err.Error(), mac, id))
	}

	cmd := []byte(fmt.Sprintf(`
			{
				"execute": "netdev_add",
				"arguments": {
					"id": "netdev-%s",
					"type": "tap",
					"ifname": "%s",
					"vhost": true,
					"script": "no",
					"downscript": "no"
				}
			}
		`, id, tap))
	raw, err := k.qmp.Run(cmd)
	if err != nil { // already existsが発行されてしまう
		return notification.MakeNotification("attachNICs.netdev_add", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	cmd = []byte(fmt.Sprintf(`
				{
					"execute": "device_add",
					"arguments": {
						"driver": "virtio-net-pci",
						"id": "virtio-net-%s",
						"netdev": "netdev-%s",
						"bus": "pci.0",
						"mac": "%s"
					}
				}
			`, id, id, hw.String()))

	raw, err = k.qmp.Run(cmd)
	if err != nil { // already existsが発行されてしまう
		return notification.MakeNotification("attachNICs.device_add", false, fmt.Sprintf("error message '%s', qmp response '%s'", err.Error(), raw))
	}

	return notification.MakeNotification("attachNICs", true, fmt.Sprintf("id:%s", id))
}

// kill $qemu
func (k kvm) kill() *pnotification.Notification {
	p, _ := os.FindProcess(k.pid)
	if err := p.Kill(); err != nil {
		return notification.MakeNotification("kill", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return notification.MakeNotification("kill", true, "")
}

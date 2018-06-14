package kvm

import (
	fmt "fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/digitalocean/go-qemu/qmp"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
)

// ps auxfww | grep $name | grep qemu
func (a KVMAgent) getProcess(name string) (*process.Process, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("Failed to get process list")
	}

	for _, p := range ps {
		c, _ := p.Cmdline()                                            // エラーが発生する場合が考えられない
		if strings.Contains(c, name) && strings.HasPrefix(c, "qemu") { // このfilterはガバガバなので後でリファクタリングする
			return p, nil
		}
	}

	return nil, nil
}

// qemu-system ...
func (a KVMAgent) startProcess(uuid uuid.UUID, name, qmpPath string, vcpus uint32, memory uint64) error {
	args := []string{
		"qemu-system-x86_64",

		// -- QEMU metadata --
		"-uuid",
		uuid.String(),
		"-name",
		fmt.Sprintf("guest=%s,debug-threads=on", name),
		"-msg",
		"timestamp=on",

		// Config
		"-daemonize",
		"-nodefaults",     // Don't create default devices
		"-no-user-config", // The "-no-user-config" option makes QEMU not load any of the user-provided config files on sysconfdir
		"-S",              // Do not start CPU at startup
		"-no-shutdown",    // Don't exit QEMU on guest shutdown

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
		"en-us", // vm.Spec.Keymapみたいなので取得できるようにする

		// VNC
		"-vnc",
		":10,websocket=5710", // TODO: ぶつからないようにポートを設定する必要がある、現状一台しか立たない

		// clock
		"-rtc",
		"base=utc,driftfix=slew",
		"-global",
		"kvm-pit.lost_tick_policy=delay",
		"-no-hpet",

		// CPU
		// TODO: 必要があればmonitorを操作してhotaddできるようにする
		// TODO: スケジューリングが可能かどうか調べる
		"-cpu",
		"host",
		"-smp",
		fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", vcpus, vcpus),
		"-enable-kvm",

		// Memory
		// TODO: スケジューリングが可能かどうか調べる
		"-m",
		fmt.Sprintf("%s", bytefmt.ByteSize(memory)),
		"-device",
		"virtio-balloon-pci,id=balloon0,bus=pci.0", // dynamic configurations
		"-realtime",
		"mlock=off",

		// VGA controller
		"-device",
		"VGA,id=video0,bus=pci.0",

		// SCSI controller
		"-device",
		"lsi53c895a,bus=pci.0,id=scsi0",
	}

	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start process, err:'%s', args:'%s'", err.Error(), args)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):
		return nil
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Failed to run process, err:'%s', args:'%s'", err.Error(), args) // stderrを表示できるようにする必要がある
		}

		return nil
	}
}

// qmp-shell .../monitor.sock
// TODO: 他のプロセスがソケットにつなげていた場合、何故か無制限にロックされてしまう
func (a KVMAgent) connectQMP(qmpPath string) (*qmp.SocketMonitor, error) {
	qmp, err := qmp.NewSocketMonitor("unix", qmpPath, 2*time.Second)
	if err != nil {
		return nil, err
	}

	return qmp, nil
}

// (QEMU) blockdev-add options={"driver":"qcow2","id":"drive-virtio-disk0","file":{"driver":"file","filename":"/home/h-otter/wk/test-qemu/ubuntu16.04.qcow2"}}
// (QEMU) device_add driver=virtio-blk-pci bus=pci.0 scsi=off drive=drive-virtio-disk0 id=virtio-disk0 bootindex=1
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - bootindexがどうやって更新されるのかがわからない
func (a KVMAgent) attachVolume(q *qmp.SocketMonitor, label string, u *url.URL, index int32) error {
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
			`, label, u.Path))
	}

	raw, err := q.Run(cmd)
	if err != nil && !strings.Contains(string(raw), "Already exists") { // TODO: contains周りの動作確認
		return fmt.Errorf("Failed to run blockdev-add, err:'%s', raw:'%s'", err.Error(), raw)
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
			`, label, label, index)) // bootindexはcdのために1を追加する

	raw, err = q.Run(cmd)
	if err != nil && !strings.Contains(string(raw), "Already exists") { // TODO: contains周りの動作確認
		return fmt.Errorf("Failed to run device_add, err:'%s', raw:'%s'", err.Error(), raw)
	}

	return nil
}

// (QEMU) netdev_add id=tap0 type=tap vhost=true ifname=tap0 script=no downscript=no
// (QEMU) device_add driver=virtio-net-pci netdev=tap0 id=test0 mac=52:54:00:df:89:29 bus=pci.0
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - MACアドレスを変更する
func (a KVMAgent) attachNIC(q *qmp.SocketMonitor, label, tap string, mac net.HardwareAddr) error {
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
		`, label, tap))
	raw, err := q.Run(cmd)
	if err != nil && !strings.Contains(string(raw), "Already exists") { // TODO: contains周りの動作確認
		return fmt.Errorf("Failed to run netdev_add, err:'%s', raw:'%s'", err.Error(), raw)
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
			`, label, label, mac.String()))

	raw, err = q.Run(cmd)
	if err != nil && !strings.Contains(string(raw), "Already exists") { // TODO: contains周りの動作確認
		return fmt.Errorf("Failed to run device_add, err:'%s', raw:'%s'", err.Error(), raw)
	}

	return nil
}

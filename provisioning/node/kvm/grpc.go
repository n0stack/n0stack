package kvm

import (
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"code.cloudfoundry.org/bytefmt"
	"github.com/digitalocean/go-qemu/qmp"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
	context "golang.org/x/net/context"
)

type KVMAgent struct{}

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
	args := []string{"qemu-system-x86_64"}

	// -- QEMU metadata --
	args = append(args, "-uuid")
	args = append(args, uuid.String())
	args = append(args, "-name")
	args = append(args, fmt.Sprintf("guest=%s,debug-threads=on", name))
	args = append(args, "-msg")
	args = append(args, "timestamp=on")

	args = append(args, "-daemonize")
	args = append(args, "-nodefaults")     // Don't create default devices
	args = append(args, "-no-user-config") // The "-no-user-config" option makes QEMU not load any of the user-provided config files on sysconfdir
	args = append(args, "-S")              // Do not start CPU at startup
	args = append(args, "-no-shutdown")    // Don't exit QEMU on guest shutdown

	// QMP
	const monitorFile = "monitor.sock"
	args = append(args, "-chardev")
	args = append(args, fmt.Sprintf("socket,id=charmonitor,path=%s,server,nowait", qmpPath))
	args = append(args, "-mon")
	args = append(args, "chardev=charmonitor,id=monitor,mode=control")

	// -- BIOS --
	// boot priority
	args = append(args, "-boot")
	args = append(args, "menu=on,strict=on")

	// keyboard
	args = append(args, "-k")
	args = append(args, "en-us") // vm.Spec.Keymapみたいなので取得できるようにする

	// VNC
	args = append(args, "-vnc")
	args = append(args, ":0,websocket=5700") // TODO: ぶつからないようにポートを設定する必要がある、現状一台しか立たない

	// clock
	args = append(args, "-rtc")
	args = append(args, "base=utc,driftfix=slew")
	args = append(args, "-global")
	args = append(args, "kvm-pit.lost_tick_policy=delay")
	args = append(args, "-no-hpet")

	// CPU
	// TODO: 必要があればmonitorを操作してhotaddできるようにする
	// TODO: スケジューリングが可能かどうか調べる
	args = append(args, "-cpu")
	args = append(args, "host")
	args = append(args, "-smp")
	args = append(args, fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", vcpus, vcpus))
	args = append(args, "-enable-kvm")
	// return true, "Succeeded to check cpu configurations"

	// Memory
	// TODO: スケジューリングが可能かどうか調べる
	args = append(args, "-m")
	args = append(args, fmt.Sprintf("%s", bytefmt.ByteSize(memory)))
	args = append(args, "-device")
	args = append(args, "virtio-balloon-pci,id=balloon0,bus=pci.0") // dynamic configurations
	args = append(args, "-realtime")
	args = append(args, "mlock=off")

	// VGA controller
	args = append(args, "-device")
	args = append(args, "VGA,id=video0,bus=pci.0")

	// SCSI controller
	args = append(args, "-device")
	args = append(args, "lsi53c895a,bus=pci.0,id=scsi0")

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
			return fmt.Errorf("Failed to run process, err:'%s', args:'%s'", err.Error(), args)
		}

		return fmt.Errorf("Failed to run process, args:'%s'", args) // stderrを表示できるようにする必要がある
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

func (a KVMAgent) ApplyKVM(ctx context.Context, req *ApplyKVMRequest) (*KVM, error) {
	// validation

	p, err := a.getProcess(req.Kvm.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to getProcess, err:'%s'", err.Error())
	}

	started := false
	if p == nil {
		u, err := uuid.FromString(req.Kvm.Uuid)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse uuid, err:'%s', uuid:'%s'", err.Error(), req.Kvm.Uuid)
		}

		err = a.startProcess(
			u,
			req.Kvm.Name,
			req.Kvm.QmpPath,
			req.Kvm.CpuCores,
			req.Kvm.MemoryBytes,
		)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to startProcess, err:'%s'", err.Error())
		}

		started = true
	}

	q, err := a.connectQMP(req.Kvm.QmpPath)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to connectQMP, err:'%s'", err.Error())
	}

	// Volume
	for label, v := range req.Kvm.Volumes {
		index := v.BootIndex
		// if index == 0 {
		// 	index = i + 10 // prefix is 10
		// }

		u, err := url.Parse(v.Url)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse url, err:'%s', url:'%s'", err.Error(), v.Url)
		}

		if err := a.attachVolume(q, label, u, index); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to attachVolume, err:'%s'", err.Error())
		}
	}

	// Network
	for label, v := range req.Kvm.Nics {

		m, err := net.ParseMAC(v.HwAddr)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse hardware address, err:'%s', hwaddr:'%s'", err.Error(), v.HwAddr)
		}

		if err := a.attachNIC(q, label, v.TapName, m); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to attachNIC, err:'%s'", err.Error())
		}
	}

	if started {
		_, err := a.Boot(context.Background(), &ActionKVMRequest{
			Name:    req.Kvm.Name,
			QmpPath: req.Kvm.QmpPath,
		})
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to Boot, err:'%s'", err.Error())
		}
	}

	return nil, nil
}

func (a KVMAgent) DeleteKVM(ctx context.Context, req *DeleteKVMRequest) (*google_protobuf.Empty, error) {
	p, err := a.getProcess(req.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to getProcess, err:'%s'", err.Error())
	}

	if err := p.Kill(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to kill process, err:'%s', pid:'%d'", err.Error(), p.Pid)
	}

	return &google_protobuf.Empty{}, nil
}

// (QEMU) cont
func (a KVMAgent) Boot(ctx context.Context, req *ActionKVMRequest) (*KVM, error) {
	res := &KVM{}

	q, err := a.connectQMP(req.QmpPath)
	if err != nil {
		return res, grpc.Errorf(codes.Internal, "Failed to connect QMP, err:'%s'", err.Error())
	}

	cmd := []byte(`{ "execute": "cont" }`)
	_, err = q.Run(cmd) // TODO: responseの結果で動作をちゃんと分ける
	if err != nil {
		res.RunLevel = KVM_SHUTDOWN // これが正しいのかはわからない

		return res, nil
	}

	res.RunLevel = KVM_RUNNING

	return res, nil
}

func (a KVMAgent) Reboot(context.Context, *ActionKVMRequest) (*KVM, error) {
	return nil, nil
}

func (a KVMAgent) HardReboot(context.Context, *ActionKVMRequest) (*KVM, error) {
	return nil, nil
}

func (a KVMAgent) Shutdown(context.Context, *ActionKVMRequest) (*KVM, error) {
	return nil, nil
}

func (a KVMAgent) HardShutdown(context.Context, *ActionKVMRequest) (*KVM, error) {
	return nil, nil
}

func (a KVMAgent) Save(context.Context, *ActionKVMRequest) (*KVM, error) {
	return nil, nil
}

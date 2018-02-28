package kvm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/target/nic/tap"
)

const NetworkType = "kvm"

type KVM struct {
	args []string
	// state string
}

func (k KVM) ManagingType() string {
	return filepath.Join(model.VMType, NetworkType)
}

// func (k KVM) Operations(state, task string) ([]string, error) {
// var err error
// if !model.NetworkStateMachine[state][task] {
// 	err = fmt.Errorf("Not allowed to operate task '%v' from state '%v'", task, state)
// }

// switch task {
// case "Up":
// 	return []string{
// 		"ParseModel",
// 		"CheckState",
// 		"CheckInterface",
// 		"CreateBridge",
// 		"UpBridge",
// 	}, err
// case "Down":
// 	return []string{
// 		"ParseModel",
// 		"CheckState",
// 		"CheckInterface",
// 		"CreateBridge",
// 		"DownBridge",
// 	}, err
// case "Delete":
// 	return []string{
// 		"ParseModel",
// 		"CheckState",
// 		"DeleteBridge",
// 	}, err
// }

// return nil, fmt.Errorf("Unsupported task '%v'", task)
// }

func (k *KVM) ParseModel(m model.AbstractModel) (bool, string) {
	v, ok := m.(*model.VM)
	if !ok {
		return false, fmt.Sprintf("Failed to parse AbstractModel to *model.VM")
	}

	return true, fmt.Sprintf("Succeeded to parse AbstractModel to *model.VM")
}

// func (k *KVM) CheckState(m model.AbstractModel) (bool, string) {
// 	n := m.(*model.Network)

// 	var err error
// 	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
// 		return false, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
// 	}

// 	f.bridgeLink.Attrs().Flags & unix.IFF_UP

// 	return true, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
// }

func (k *KVM) CheckCPU(m model.AbstractModel) (bool, string) {
	v := m.(*model.VM)

	// TODO: 必要があればmonitorを操作してhotaddできるようにする
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-cpu")
	k.args = append(k.args, "host")
	k.args = append(k.args, "-smp")
	k.args = append(k.args, fmt.Sprintf("%d,sockets=%d,cores=1,threads=1", v.VCPUs, v.VCPUs))

	return true, "Succeeded to check cpu configurations"
}

// -m 512M
func (k *KVM) CheckMemory(m model.AbstractModel) (bool, string) {
	vm := m.(*model.VM)

	// TODO: 必要があればmonitorを操作してhotaddできるようにする
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-m")
	k.args = append(k.args, fmt.Sprintf("%v", vm.Memory)) // メモリの単位を検証する

	return true, "Succeeded to check memory configurations"
}

// -boot order=c -drive file=ubuntu16.04.qcow2,format=qcow2,if=virtio,index=0
func (k *KVM) CheckVolume(m model.AbstractModel) (bool, string) {
	vm := m.(*model.VM)

	bootOnce := ""
	bootOrder := ""

	for _, d := range vm.Dependencies.SelectWithModelType("resource/volume") {
		v := d.Model.(*model.Volume)

		switch v.Type {
		case "resource/volume/iso":
			k.args = append(k.args, "-cdrom")
			k.args = append(k.args, fmt.Sprintf("%v", v.URL))

		case "resource/volume/local_qcow2":
			k.args = append(k.args, "-drive")
			k.args = append(k.args, fmt.Sprintf("file=%v,format=qcow2,if=virtio,index=0", v.URL)) // indexについて
		}
	}

	k.args = append(k.args, "-boot")
	k.args = append(k.args, "order=c")
	k.args = append(k.args, "order=c,once=d")

	return true, ""
}

// -net nic,macaddr=52:54:11:11:11:11,model=virtio -net tap,ifname=hoge,script=no,downscript=no
func (k *KVM) CheckInterface(m model.AbstractModel) (bool, string) {
	vm := m.(*model.VM)

	for _, d := range vm.Dependencies.SelectWithLabel("n0stack/n0core/resource/vm/nic") {
		nic := d.Model.(*model.NIC)
		switch nic.Type {
		case filepath.Join(model.NICType, tap.NICType):
			k.args = append(k.args, "-net")
			k.args = append(k.args, fmt.Sprintf("nic,macaddr=%v,model=virtio", nic.HWAddr.String()))
			k.args = append(k.args, "-net")
			k.args = append(k.args, fmt.Sprintf("nic,tap,ifname=%v,script=no,downscript=no", nic.Meta["tap_name"]))
		}
	}

	return true, ""
}

// sudo qemu-system-x86_64 -cpu host -enable-kvm -boot order=c -drive file=ubuntu16.04.qcow2,format=qcow2,if=virtio,index=0 -smp 2,sockets=2,cores=1,threads=1 -m 512M -net nic,macaddr=52:54:11:11:11:11,model=virtio -net tap,ifname=hoge,script=no,downscript=no
func (k *KVM) LaunchVM(m model.AbstractModel) (bool, string) {
	vm := m.(*model.VM)

	k.args = append(k.args, "-enable-kvm")
	k.args = append(k.args, "-qmp")
	k.args = append(k.args, fmt.Sprintf("unix:/tmp/n0core-qmp%v.sock,server,nowait", vm.ID.String()))
	k.args = append(k.args, "-uuid")
	k.args = append(k.args, vm.ID.String())

	qemu := fmt.Sprintf("qemu-system-%s", vm.Arch)
	cmd := exec.Command(qemu, k.args...)

	err := cmd.Start()
	if err != nil {
		return false, fmt.Sprintf("Failed to start qemu process: error message '%v'", err.Error())
	}

	cmd.Wait()

	return true, fmt.Sprintf("Succeeded to launch vm '%v'", vm.ID)
}

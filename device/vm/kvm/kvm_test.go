package kvm

import (
	"strings"
	"testing"
)

func TestGetQMPPath(t *testing.T) {
	a := strings.Split("qemu-system-x86_64 -uuid 614de68e-8e7c-429e-830f-76593a1d69ce -name guest=n0core-test-vm,debug-threads=on -msg timestamp=on -nodefaults -no-user-config -S -no-shutdown -chardev socket,id=charmonitor,path=/var/lib/n0core/device/vm/kvm/614de68e-8e7c-429e-830f-76593a1d69ce/monitor.sock,server,nowait -mon chardev=charmonitor,id=monitor,mode=control -boot menu=on,strict=on -k en-us -vnc :0 -rtc base=utc,driftfix=slew -global kvm-pit.lost_tick_policy=delay -no-hpet -cpu host -smp 1,sockets=1,cores=1,threads=1 -enable-kvm -m 4G -device virtio-balloon-pci,id=balloon0,bus=pci.0 -realtime mlock=off -device VGA,id=video0,bus=pci.0 -device lsi53c895a,bus=pci.0,id=scsi0", " ")
	k := kvm{
		args: a,
	}

	if k.getQMPPath() != "/var/lib/n0core/device/vm/kvm/614de68e-8e7c-429e-830f-76593a1d69ce/monitor.sock" {
		t.Errorf("got %v\nwant %v", k.getQMPPath(), "/var/lib/n0core/device/vm/kvm/614de68e-8e7c-429e-830f-76593a1d69ce/monitor.sock")
	}
}

package node

import (
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// IPMIを持っていない場合が考えられるので、とりあえずエラーハンドリングはしていない
func GetIpmiAddress() string {
	out, err := exec.Command("ipmitool", "lan", "print").Output()
	if err != nil {
		return ""
	}

	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "IP Address              :") { // これで正しいのかよくわからず、要テスト
			s := strings.Split(l, " ")
			return s[len(s)-1]
		}
	}

	return ""
}

// Serialが取得できなくても動作に問題はないため、エラーハンドリングはしていない
func GetSerial() string {
	out, err := exec.Command("dmidecode", "-t", "system").Output()
	if err != nil {
		return ""
	}

	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "Serial Number:") {
			s := strings.Split(l, " ")
			return s[len(s)-1]
		}
	}

	return ""
}

func GetTotalMemory() uint64 {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	if err != nil {
		return 0
	}

	return si.Totalram / uint64(si.Unit)
}

func GetTotalCPUMilliCores() uint32 {
	return uint32(runtime.NumCPU())
}

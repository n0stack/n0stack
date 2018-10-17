package dnsmasq

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Dnsmasq struct {
	pid  int32
	proc *process.Process

	baseDir string

	hostEntries map[string]string
}

func OpenDnsmasq(baseDir string) (*Dnsmasq, error) {
	d := &Dnsmasq{
		baseDir: baseDir,
	}

	d.getPid()
	d.readConfigFiles()

	return d, nil
}

func (d *Dnsmasq) IsRunning() bool {
	if ok, _ := process.PidExists(d.pid); !ok {
		return false
	}

	return true
}

func (d Dnsmasq) getPid() error {
	return nil
}
func (d Dnsmasq) pidfile() string {
	return filepath.Join(d.baseDir, "pid")
}
func (d Dnsmasq) hostsfile() string {
	return filepath.Join(d.baseDir, "hostsfile")
}
func (d Dnsmasq) leasefile() string {
	return filepath.Join(d.baseDir, "pid")
}

func (d *Dnsmasq) Start(interfaceName string) (io.Writer, io.Writer, error) {
	args := []string{
		"/usr/sbin/dnsmasq",
		"--no-hosts",
		"--no-resolv",
		"--keep-in-foreground",
		// "--bind-interfaces",
		fmt.Sprintf("--interface=%s", interfaceName),
		"--except-interface=lo",
		fmt.Sprintf("--pid-file=%s", d.pidfile()),
		fmt.Sprintf("--dhcp-hostsfile=%s", d.hostsfile()),
		fmt.Sprintf("--dhcp-leasefile=%s", d.leasefile()),
		// fmt.Sprintf("--dhcp-optsfile=%s", d.optsfile()),
		// fmt.sprintf("--dhcp-range=%s,%s,12h"),
		// fmt.sprintf("--log-facility=%s"),
	}

	cmd := exec.Command(args[0], args[1:]...)
	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("Failed to start process: args='%s', err='%s'", args, err.Error())
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(1 * time.Second):
		break

	case err := <-done:
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to run process: args='%s', err='%s'", args, err.Error()) // stderrを表示できるようにする必要がある
		}
	}

	return cmd.Stdout, cmd.Stderr, nil
}

func (d *Dnsmasq) Delete() error {
	return nil
}

func (d *Dnsmasq) readConfigFiles() error {
	return nil
}

func (d *Dnsmasq) writeHostfile() error {
	f, err := os.OpenFile(d.hostsfile(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open hostsfile: hostsfile='%s', err='%s'", d.hostsfile(), err.Error())
	}
	defer f.Close()

	for k, v := range d.hostEntries {
		fmt.Fprintf(f, "%s,%s\n", k, v)
	}

	return nil
}

func (d *Dnsmasq) applyConfigFiles() error {
	if err := d.writeHostfile(); err != nil {
		return err
	}

	if !d.IsRunning() {
		return fmt.Errorf("")
	}

	if err := d.proc.SendSignal(syscall.SIGHUP); err != nil {
		return fmt.Errorf("Failed to reload config: err='%s'", err.Error())
	}

	return nil
}

func (d *Dnsmasq) ApplyHostEntry(mac net.HardwareAddr, ip net.IP) error {
	d.hostEntries[mac.String()] = ip.String()

	if err := d.applyConfigFiles(); err != nil {
		return fmt.Errorf("Failed to apply config: err='%s'", err.Error())
	}

	return nil
}

func (d *Dnsmasq) DeleteHostEntry(mac net.HardwareAddr) error {
	if _, ok := d.hostEntries[mac.String()]; ok {
		delete(d.hostEntries, mac.String())
	} else {
		// return fmt.Errorf("")
	}

	if err := d.applyConfigFiles(); err != nil {
		return fmt.Errorf("Failed to apply config: err='%s'", err.Error())
	}

	return nil
}

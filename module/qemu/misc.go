package qemu

import (
	"fmt"
	"strings"
	"time"

	"github.com/digitalocean/go-qemu/qmp"
	"github.com/digitalocean/go-qemu/qmp/raw"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
)

type Qemu struct {
	proc *process.Process

	// args
	id      *uuid.UUID
	qmpPath string

	qmp qmp.Monitor
	m   *raw.Monitor
}

func (q *Qemu) init() error {
	if err := q.findProcess(q.id.String()); err != nil {
		return fmt.Errorf("Failed to find process: err='%s'", err.Error())
	}

	if q.proc != nil {
		c, err := q.proc.CmdlineSlice()
		if err != nil {
			return fmt.Errorf("Failed to get command line: err='%s'", err.Error())
		}

		if err := q.parseArgs(c); err != nil {
			return fmt.Errorf("Failed to parse arguments of command line: err='%s'", err.Error())
		}

		if err := q.initQMP(); err != nil {
			return fmt.Errorf("Failed to initialize QMP socket: err='%s'", err.Error())
		}
	}

	return nil
}

func (q *Qemu) findProcess(contain string) error {
	ps, err := process.Processes()
	if err != nil {
		return fmt.Errorf("Failed to get process list")
	}

	for _, p := range ps {
		c, _ := p.Cmdline()                                               // エラーが発生する場合が考えられない
		if strings.Contains(c, contain) && strings.HasPrefix(c, "qemu") { // このfilterはガバガバなので後でリファクタリングする
			q.proc = p

			return nil
		}
	}

	return nil
}

func (q *Qemu) initQMP() error {
	qmp, err := qmp.NewSocketMonitor("unix", q.qmpPath, 2*time.Second)
	if err != nil {
		return fmt.Errorf("Failed to open QMP socket: err='%s'", err.Error())
	}

	if err := qmp.Connect(); err != nil {
		return fmt.Errorf("Failed to connect QMP socket: err='%s'", err.Error())
	}

	q.qmp = qmp
	q.m = raw.NewMonitor(q.qmp)

	return nil
}

func (q *Qemu) parseArgs(args []string) error {
	q.qmpPath = "monitor.sock"

	return nil
}

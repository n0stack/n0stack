package qemu

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/digitalocean/go-qemu/qmp"
	"github.com/digitalocean/go-qemu/qmp/raw"
	"github.com/shirou/gopsutil/process"
)

func (q *Qemu) init() error {
	if err := q.findProcess(q.name); err != nil {
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
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5) // CIがこけてしまうため長くしているが、かなり長時間待つようになっているため普通の処理においてタイムアウトする危険性がある
	err := backoff.Retry(func() (err error) {
		q.qmp, err = qmp.NewSocketMonitor("unix", q.qmpPath, 3*time.Second)
		return
	}, b)
	if err != nil {
		return fmt.Errorf("Failed to open QMP socket: err='%s'", err.Error())
	}

	if err := q.qmp.Connect(); err != nil {
		return fmt.Errorf("Failed to connect QMP socket: err='%s'", err.Error())
	}

	q.m = raw.NewMonitor(q.qmp)

	return nil
}

func (q *Qemu) parseArgs(args []string) error {
	mon := GetQemuArgValue("-mon", "*", args)
	if mon != nil {
		_, ok := mon.kwds["chardev"]
		if ok {
			mc := GetQemuArgValue("-chardev", mon.kwds["chardev"], args)
			q.qmpPath = mc.kwds["path"]
		}
	}

	return nil
}

type QemuArgValue struct {
	args []string
	kwds map[string]string
}

// GetQemuArgValue is simple qemu argument parser.
// Order is O(len(args)).
//
// Example:
//   `-mon chardev=charmonitor,id=monitor`
//
//   Args: option="-mon", id="monitor"
//         "*" is wild card, return only first match
//   Retrun: {"arg": "", "kwds": {"chardev": "charmonitor", "id": "monitor"}}
func GetQemuArgValue(option, id string, args []string) *QemuArgValue {
	q := &QemuArgValue{
		kwds: map[string]string{},
	}

	for i, a := range args {
		if option == "*" || option == a {
			values := strings.Split(args[i+1], ",")

			for _, v := range values {
				kv := strings.Split(v, "=")
				if len(kv) == 1 {
					q.args = append(q.args, v)
					continue
				}

				q.kwds[kv[0]] = kv[1]
			}

			if id == "*" || q.kwds["id"] == id {
				return q
			}
		}
	}

	return nil
}

func GetNewListenPort(begin uint16) uint16 {
	for p := begin; p <= uint16(65535); p++ {
		l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", p))
		if err == nil {
			defer l.Close()

			return p
		}
	}

	return 0
}

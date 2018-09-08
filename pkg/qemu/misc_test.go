package qemu

import (
	"reflect"
	"testing"
)

func TestGetQemuArgValue(t *testing.T) {
	Cases := []struct {
		name   string
		option string
		id     string
		args   []string
		result *QemuArgValue
	}{
		{
			"with arg",
			"-chardev",
			"charmonitor",
			[]string{"-chardev", "socket,id=charmonitor,path=monitor.sock,server,nowait"},
			&QemuArgValue{
				args: []string{"socket", "server", "nowait"},
				kwds: map[string]string{
					"id":   "charmonitor",
					"path": "monitor.sock",
				},
			},
		},
		{
			"without arg",
			"-mon",
			"monitor",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			&QemuArgValue{
				kwds: map[string]string{
					"id":      "monitor",
					"chardev": "charmonitor",
					"mode":    "control",
				},
			},
		},
		{
			"wild card option",
			"*",
			"monitor",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			&QemuArgValue{
				kwds: map[string]string{
					"id":      "monitor",
					"chardev": "charmonitor",
					"mode":    "control",
				},
			},
		},
		{
			"wild card id",
			"-mon",
			"*",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			&QemuArgValue{
				kwds: map[string]string{
					"id":      "monitor",
					"chardev": "charmonitor",
					"mode":    "control",
				},
			},
		},
		{
			"no match",
			"-foo",
			"*",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			nil,
		},
	}

	for _, c := range Cases {
		v := GetQemuArgValue(c.option, c.id, c.args)

		if !reflect.DeepEqual(v, c.result) {
			t.Errorf("[%s] Wrong result\n\thave:%v\n\twant:%v", c.name, v, c.result)
		}
	}
}

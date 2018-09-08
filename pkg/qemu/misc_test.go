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
		result *qemuArgValue
	}{
		{
			"with arg",
			"-chardev",
			"charmonitor",
			[]string{"-chardev", "socket,id=charmonitor,path=monitor.sock,server,nowait"},
			&qemuArgValue{
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
			&qemuArgValue{
				kwds: map[string]string{
					"id":      "monitor",
					"chardev": "charmonitor",
					"mode":    "control",
				},
			},
		},
	}

	//-mon chardev=charmonitor,id=monitor
	for _, c := range Cases {
		v := getQemuArgValue(c.option, c.id, c.args)

		if v == nil {
			t.Errorf("[%s] Result is nil\n\twant:%v", c.name, c.result)
		}

		if !reflect.DeepEqual(v, c.result) {
			t.Errorf("[%s] Wrong result\n\thave:%v\n\twant:%v", c.name, v, c.result)
		}
	}
}

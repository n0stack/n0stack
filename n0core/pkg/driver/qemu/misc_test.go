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

func TestGetParsedOptionValueById(t *testing.T) {
	Cases := []struct {
		name string
		key  string
		id   string
		req  []string
		args []string
		kwds map[string]string
		ok   bool
	}{
		{
			"with arg",
			"-chardev",
			"charmonitor",
			[]string{"-chardev", "socket,id=charmonitor,path=monitor.sock,server,nowait"},
			[]string{"socket", "server", "nowait"},
			map[string]string{
				"id":   "charmonitor",
				"path": "monitor.sock",
			},
			true,
		},
		{
			"without arg",
			"-mon",
			"monitor",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			[]string{},
			map[string]string{
				"id":      "monitor",
				"chardev": "charmonitor",
				"mode":    "control",
			},
			true,
		},
		{
			"no match",
			"-foo",
			"aa",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			nil,
			nil,
			false,
		},
	}

	for _, c := range Cases {
		a := ParseQemuArgs(c.req)

		args, kwds, ok := a.GetParsedOptionValueById(c.key, c.id)

		if c.ok != ok {
			t.Errorf("[%s] ok is wrong", c.name)
		}
		if diff := cmp.Diff(c.args, args); diff != "" {
			t.Errorf("[%s] args is wrong: diff=(-want +got)\n%s", c.name, diff)
		}
		if diff := cmp.Diff(c.kwds, kwds); diff != "" {
			t.Errorf("[%s] kwds is wrong: diff=(-want +got)\n%s", c.name, diff)
		}
	}
}

func TestGetTopParsedOptionValue(t *testing.T) {
	Cases := []struct {
		name string
		key  string
		req  []string
		args []string
		kwds map[string]string
		ok   bool
	}{
		{
			"with arg",
			"-chardev",
			[]string{"-chardev", "socket,id=charmonitor,path=monitor.sock,server,nowait"},
			[]string{"socket", "server", "nowait"},
			map[string]string{
				"id":   "charmonitor",
				"path": "monitor.sock",
			},
			true,
		},
		{
			"without arg",
			"-mon",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			[]string{},
			map[string]string{
				"id":      "monitor",
				"chardev": "charmonitor",
				"mode":    "control",
			},
			true,
		},
		{
			"no match",
			"-foo",
			[]string{"-mon", "chardev=charmonitor,id=monitor,mode=control"},
			nil,
			nil,
			false,
		},
	}

	for _, c := range Cases {
		a := ParseQemuArgs(c.req)

		args, kwds, ok := a.GetTopParsedOptionValue(c.key)

		if c.ok != ok {
			t.Errorf("[%s] ok is wrong", c.name)
		}
		if diff := cmp.Diff(c.args, args); diff != "" {
			t.Errorf("[%s] args is wrong: diff=(-want +got)\n%s", c.name, diff)
		}
		if diff := cmp.Diff(c.kwds, kwds); diff != "" {
			t.Errorf("[%s] kwds is wrong: diff=(-want +got)\n%s", c.name, diff)
		}
	}
}

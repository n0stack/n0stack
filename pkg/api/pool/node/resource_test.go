package node

import (
	"testing"

	"github.com/n0stack/proto.go/resource/v0"
)

func TestCheckCompute(t *testing.T) {
	cases := []struct {
		name     string
		cpu      uint32
		mem      uint64
		total    *presource.Compute
		reserved map[string]*presource.Compute
		err      string
	}{
		{
			name: "Valid: no reserve, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1 * 10,
				MemoryBytes: 1024 * 10,
			},
			reserved: map[string]*presource.Compute{},
			err:      "",
		},
		{
			name: "Valid: no reserve, cpu_total=cpu_request, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1,
				MemoryBytes: 1024 * 10,
			},
			reserved: map[string]*presource.Compute{},
			err:      "",
		},
		{
			name: "Valid: no reserve, mem_total=mem_request, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1 * 10,
				MemoryBytes: 1024,
			},
			reserved: map[string]*presource.Compute{},
			err:      "",
		},
		{
			name: "Valid: some reserve, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1 * 10,
				MemoryBytes: 1024 * 10,
			},
			reserved: map[string]*presource.Compute{
				"hoge": &presource.Compute{
					Vcpus: 1,
				},
			},
			err: "",
		},
		{
			name: "Invalid: no reserve, over cpu",
			cpu:  2,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1,
				MemoryBytes: 1024 * 10,
			},
			reserved: map[string]*presource.Compute{},
			err:      "parameter='vcpus', total='1', used='0', requested='2'",
		},
		{
			name: "Invalid: no reserve, over memory",
			cpu:  1,
			mem:  1025,
			total: &presource.Compute{
				Vcpus:       1,
				MemoryBytes: 1024,
			},
			reserved: map[string]*presource.Compute{},
			err:      "parameter='memory', total='1024', used='0', requested='1025'",
		},
		{
			name: "Invalid: no reserve, over cpu, memory",
			cpu:  2,
			mem:  1025,
			total: &presource.Compute{
				Vcpus:       1,
				MemoryBytes: 1024,
			},
			reserved: map[string]*presource.Compute{},
			err:      "parameter='vcpus', total='1', used='0', requested='2'",
		},
		{
			name: "Invalid: some reserve, over cpu",
			cpu:  1,
			mem:  1024,
			total: &presource.Compute{
				Vcpus:       1,
				MemoryBytes: 1024,
			},
			reserved: map[string]*presource.Compute{
				"hoge": &presource.Compute{
					Vcpus: 1,
				},
			},
			err: "parameter='vcpus', total='1', used='1', requested='1'",
		},
	}

	for _, c := range cases {
		err := CheckCompute(c.cpu, c.mem, c.total, c.reserved)

		if err == nil && c.err != "" {
			t.Errorf("[%s] Need error message\n\thave:nil\n\twant:%s", c.name, c.err)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("[%s] Wrong error message\n\thave:%s\n\twant:%s", c.name, err.Error(), c.err)
		}
	}
}

func TestCheckStorage(t *testing.T) {
	cases := []struct {
		name     string
		req      uint64
		total    *presource.Storage
		reserved map[string]*presource.Storage
		err      string
	}{
		{
			name: "Valid: no reserve, total=request*10",
			req:  1024,
			total: &presource.Storage{
				Bytes: 1024 * 10,
			},
			reserved: map[string]*presource.Storage{},
			err:      "",
		},
		{
			name: "Valid: no reserve, total=req",
			req:  1024,
			total: &presource.Storage{
				Bytes: 1024,
			},
			reserved: map[string]*presource.Storage{},
			err:      "",
		},
		{
			name: "Valid: some reserve, total=request*10",
			req:  1024,
			total: &presource.Storage{
				Bytes: 1024 * 10,
			},
			reserved: map[string]*presource.Storage{
				"hoge": &presource.Storage{
					Bytes: 1024,
				},
			},
			err: "",
		},
		{
			name: "Invalid: no reserve, over req",
			req:  1025,
			total: &presource.Storage{
				Bytes: 1024,
			},
			reserved: map[string]*presource.Storage{},
			err:      "total='1024', used='0', requested='1025'",
		},
		{
			name: "Invalid: some reserve, over req",
			req:  1024,
			total: &presource.Storage{
				Bytes: 1024,
			},
			reserved: map[string]*presource.Storage{
				"hoge": &presource.Storage{
					Bytes: 1024,
				},
			},
			err: "total='1024', used='1024', requested='1024'",
		},
	}

	for _, c := range cases {
		err := CheckStorage(c.req, c.total, c.reserved)

		if err == nil && c.err != "" {
			t.Errorf("[%s] Need error message\n\thave:nil\n\twant:%s", c.name, c.err)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("[%s] Wrong error message\n\thave:%s\n\twant:%s", c.name, err.Error(), c.err)
		}
	}
}

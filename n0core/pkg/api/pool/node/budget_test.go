package node

import (
	"testing"

	"n0st.ac/n0stack/n0proto.go/budget/v0"
)

func TestCheckCompute(t *testing.T) {
	cases := []struct {
		name     string
		cpu      uint32
		mem      uint64
		total    *pbudget.Compute
		reserved map[string]*pbudget.Compute
		err      string
	}{
		{
			name: "Valid: no reserve, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1 * 10,
				RequestMemoryBytes:  1024 * 10,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "",
		},
		{
			name: "Valid: no reserve, cpu_total=cpu_request, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1,
				RequestMemoryBytes:  1024 * 10,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "",
		},
		{
			name: "Valid: no reserve, mem_total=mem_request, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1 * 10,
				RequestMemoryBytes:  1024,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "",
		},
		{
			name: "Valid: some reserve, total=request*10",
			cpu:  1,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1 * 10,
				RequestMemoryBytes:  1024 * 10,
			},
			reserved: map[string]*pbudget.Compute{
				"hoge": {
					RequestCpuMilliCore: 1,
				},
			},
			err: "",
		},
		{
			name: "Invalid: no reserve, over cpu",
			cpu:  2,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1,
				RequestMemoryBytes:  1024 * 10,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "parameter='cpu_milli_core', total='1', used='0', requested='2'",
		},
		{
			name: "Invalid: no reserve, over memory",
			cpu:  1,
			mem:  1025,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1,
				RequestMemoryBytes:  1024,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "parameter='memory_bytes', total='1024', used='0', requested='1025'",
		},
		{
			name: "Invalid: no reserve, over cpu, memory",
			cpu:  2,
			mem:  1025,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1,
				RequestMemoryBytes:  1024,
			},
			reserved: map[string]*pbudget.Compute{},
			err:      "parameter='cpu_milli_core', total='1', used='0', requested='2'",
		},
		{
			name: "Invalid: some reserve, over cpu",
			cpu:  1,
			mem:  1024,
			total: &pbudget.Compute{
				RequestCpuMilliCore: 1,
				RequestMemoryBytes:  1024,
			},
			reserved: map[string]*pbudget.Compute{
				"hoge": {
					RequestCpuMilliCore: 1,
				},
			},
			err: "parameter='cpu_milli_core', total='1', used='1', requested='1'",
		},
	}

	for _, c := range cases {
		err := CheckCompute(c.cpu, c.total.RequestCpuMilliCore, c.mem, c.total.RequestMemoryBytes, c.reserved)

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
		total    uint64
		reserved map[string]*pbudget.Storage
		err      string
	}{
		{
			name:     "Valid: no reserve, total=request*10",
			req:      1024,
			total:    1024 * 10,
			reserved: map[string]*pbudget.Storage{},
			err:      "",
		},
		{
			name:     "Valid: no reserve, total=req",
			req:      1024,
			total:    1024,
			reserved: map[string]*pbudget.Storage{},
			err:      "",
		},
		{
			name:  "Valid: some reserve, total=request*10",
			req:   1024,
			total: 1024 * 10,
			reserved: map[string]*pbudget.Storage{
				"hoge": {
					RequestBytes: 1024,
				},
			},
			err: "",
		},
		{
			name:     "Invalid: no reserve, over req",
			req:      1025,
			total:    1024,
			reserved: map[string]*pbudget.Storage{},
			err:      "total='1024', used='0', requested='1025'",
		},
		{
			name:  "Invalid: some reserve, over req",
			req:   1024,
			total: 1024,
			reserved: map[string]*pbudget.Storage{
				"hoge": {
					RequestBytes: 1024,
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

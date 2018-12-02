package netutil

import "testing"

func TestStructLinuxNetdevName(t *testing.T) {
	cases := []struct {
		name string
		arg  string
		ret  string
	}{
		{
			"simple",
			"aa",
			"aa-19d7",
		},
		{
			"trim",
			"aaaaaaaaaaaaaaaa",
			"aaaaaaaaaa-68d5",
		},
		{
			"just",
			"aaaaaaaaaa",
			"aaaaaaaaaa-cdf0",
		},
		{
			"just+1",
			"aaaaaaaaaaa",
			"aaaaaaaaaa-5d92",
		},
	}

	for _, c := range cases {
		ret := StructLinuxNetdevName(c.arg)
		if ret != c.ret {
			t.Errorf("Get wrong return, have='%s', want='%s'", ret, c.ret)
		}
	}
}

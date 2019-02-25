package stringutil

import "testing"

func TestStringWithChecksumSuffix(t *testing.T) {
	cases := []struct {
		name string
		arg  string
		size int
		ret  string
	}{
		{
			"simple",
			"aa",
			15,
			"aa-19d7",
		},
		{
			"trim",
			"aaaaaaaaaaaaaaaa",
			15,
			"aaaaaaaaaa-68d5",
		},
		{
			"just",
			"aaaaaaaaaa",
			15,
			"aaaaaaaaaa-cdf0",
		},
		{
			"just+1",
			"aaaaaaaaaaa",
			15,
			"aaaaaaaaaa-5d92",
		},
		{
			"just+1",
			"aaaaaaaaaaa",
			10,
			"aaaaa-5d92",
		},
	}

	for _, c := range cases {
		ret := StringWithChecksumSuffix(c.arg, c.size)
		if ret != c.ret {
			t.Errorf("Get wrong return, have='%s', want='%s'", ret, c.ret)
		}
	}
}

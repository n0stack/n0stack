package provisioning

import "testing"

func TestTrimNetworkInterfaceName(t *testing.T) {
	origin := "0123456789abcdefxxxx"
	want := "0123456789abcdef"
	if s := TrimNetdevName(origin); s != want {
		t.Errorf("Result is wrong: want='%s', have='%s'", want, s)
	}
}

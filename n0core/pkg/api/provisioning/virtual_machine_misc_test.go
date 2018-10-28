package provisioning

import "testing"

func TestTrimNetworkInterfaceName(t *testing.T) {
	origin := "0123456789abcdefxxxx"
	want := "0123456789abcde"
	if s := TrimNetdevName(origin); s != want {
		t.Errorf("Result is wrong: want='%s', have='%s'", want, s)
	}

	origin = "0123456789abcde"
	want = "0123456789abcde"
	if s := TrimNetdevName(origin); s != want {
		t.Errorf("Result is wrong: want='%s', have='%s'", want, s)
	}

	origin = "0123456789abcdef"
	want = "0123456789abcde"
	if s := TrimNetdevName(origin); s != want {
		t.Errorf("Result is wrong: want='%s', have='%s'", want, s)
	}

	origin = "012"
	want = "012"
	if s := TrimNetdevName(origin); s != want {
		t.Errorf("Result is wrong: want='%s', have='%s'", want, s)
	}
}

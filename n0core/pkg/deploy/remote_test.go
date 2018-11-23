package deploy

import "testing"

func TestReadSelf(t *testing.T) {
	d := &RemoteDeployer{
		targetDirectory: "/var/lib/n0core",
	}

	if buf, err := d.ReadSelf(); err != nil {
		t.Errorf("Failed to read self: err=%s'", err.Error())
	} else if len(buf) < 1 {
		t.Errorf("Length of buffer is not large, maybe wrong: len=%d'", len(buf))
	}
}

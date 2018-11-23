package deploy

import (
	"os"
	"testing"
)

func TestLocalDeployer(t *testing.T) {
	d, err := NewLocalDeployer(".")
	if err != nil {
		t.Fatalf("Failed to create new local deployer ")
	}

	if out, err := d.InstallPackages(); err != nil {
		t.Errorf("Failed to InstallPackages: err='%s'", err.Error())
	} else {
		t.Logf("The output of InstallPackages: out=\n%s", out)
	}

	testPath := "test"
	if err := d.SaveFile([]byte("test"), testPath, 0644); err != nil {
		t.Errorf("Failed to SaveFile: err='%s'", err.Error())
	} else {
		if err := os.Remove(testPath); err != nil {
			t.Fatalf("Failed to Remove test environment: err='%s'", err.Error())
		}
	}
}

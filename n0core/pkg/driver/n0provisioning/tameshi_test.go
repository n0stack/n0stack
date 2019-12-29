package n0provisioning

import (
	"testing"

	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
)

func TestNewCreateVirtualMachine(t *testing.T) {
	u := CreateVirtualMachine(&pprovisioning.CreateVirtualMachineRequest{
		Name: "testing",
		Labels: map[string]string{
			"foo": "bar",
		},
	})

	t.Errorf("%s", u.String())
}

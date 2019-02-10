package network

import (
	"testing"

	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

func TestIsLockedForDeletion(t *testing.T) {
	cases := []struct {
		Name    string
		Network *ppool.Network
		Res     bool
	}{
		{
			Name:    "nothing",
			Network: &ppool.Network{},
			Res:     false,
		},
		{
			Name: "not lock",
			Network: &ppool.Network{
				ReservedNetworkInterfaces: map[string]*pbudget.NetworkInterface{
					"not lock1": &pbudget.NetworkInterface{
						Annotations: map[string]string{
							AnnotationNetworkInterfaceDisableDeletionLock: "true",
						},
					},
					"not lock2": &pbudget.NetworkInterface{
						Annotations: map[string]string{
							AnnotationNetworkInterfaceDisableDeletionLock: "true",
						},
					},
				},
			},
			Res: false,
		},
		{
			Name: "lock",
			Network: &ppool.Network{
				ReservedNetworkInterfaces: map[string]*pbudget.NetworkInterface{
					"lock": &pbudget.NetworkInterface{},
				},
			},
			Res: true,
		},

		{
			Name: "lock after not lock",
			Network: &ppool.Network{
				ReservedNetworkInterfaces: map[string]*pbudget.NetworkInterface{
					"not lock": &pbudget.NetworkInterface{
						Annotations: map[string]string{
							AnnotationNetworkInterfaceDisableDeletionLock: "true",
						},
					},
					"lock": &pbudget.NetworkInterface{},
				},
			},
			Res: true,
		},
	}

	for _, c := range cases {
		if IsLockedForDeletion(c.Network) != c.Res {
			t.Errorf("[%s] IsLockedForDeletion return value was wrong: network=%v", c.Name, c.Network)
		}
	}
}

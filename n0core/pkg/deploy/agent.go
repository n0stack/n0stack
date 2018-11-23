package deploy

import (
	"bytes"
	"fmt"

	"github.com/coreos/go-systemd/unit"
)

func (d LocalDeployer) CreateAgentUnit(self, args string) []byte {
	u := []*unit.UnitOption{
		{
			Section: "Unit",
			Name:    "Description",
			Value:   "n0core agent: The n0stack cluster node",
		},
		{
			Section: "Unit",
			Name:    "Documentation",
			Value:   "https://github.com/n0stack/n0stack",
		},
		// {
		// 	Section: "Service",
		// 	Name:    "Environment",
		// 	Value:   "",
		// },
		{
			Section: "Service",
			Name:    "ExecStart",
			Value:   fmt.Sprintf("%s agent %s", self, args),
		},
		{
			Section: "Service",
			Name:    "Restart",
			Value:   "always",
		},
		{
			Section: "Service",
			Name:    "StartLimitInterval",
			Value:   "0",
		},
		{
			Section: "Service",
			Name:    "RestartSec",
			Value:   "10",
		},
		{
			Section: "Install",
			Name:    "WantedBy",
			Value:   "multi-user.target",
		},
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(unit.Serialize(u))
	return buf.Bytes()
}

package deploy

import (
	"bytes"

	"github.com/coreos/go-systemd/unit"
)

func (d LocalDeployer) CreateAgentUnit(command string) []byte {
	u := []*unit.UnitOption{
		{
			Section: "Unit",
			Name:    "Description",
			Value:   "n0core agent: The n0stack cluster node",
		},
		{
			Section: "Unit",
			Name:    "Documentation",
			Value:   "https://n0st.ac/n0stack",
		},
		// {
		// 	Section: "Service",
		// 	Name:    "Environment",
		// 	Value:   "",
		// },
		{
			Section: "Service",
			Name:    "ExecStart",
			Value:   command,
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
			Section: "Service",
			Name:    "KillMode",
			Value:   "process",
		},
		{
			Section: "Service",
			Name:    "TasksMax",
			Value:   "infinity",
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

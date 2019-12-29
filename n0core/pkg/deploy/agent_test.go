package deploy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateAgentUnit(t *testing.T) {
	d := &LocalDeployer{
		targetDirectory: "/var/lib/n0core",
	}

	have := string(d.CreateAgentUnit("/var/lib/n0core/n0core agent hogehoge"))
	want := `[Unit]
Description=n0core agent: The n0stack cluster node
Documentation=https://n0st.ac/n0stack

[Service]
ExecStart=/var/lib/n0core/n0core agent hogehoge
Restart=always
StartLimitInterval=0
RestartSec=10
KillMode=process
TasksMax=infinity

[Install]
WantedBy=multi-user.target
`

	if diff := cmp.Diff(have, want); diff != "" {
		t.Errorf("CreateAgentUnit response is wrong: diff=(-want +have)\n%s", diff)
	}
}

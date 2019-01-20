package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const systemdAgentUnitPath = "/etc/systemd/system/n0core-agent.service"
const systemdAgentUnit = "n0core-agent"

var AgentRequiredPackages = []string{
	"cloud-image-utils",
	"ipmitool",
	"iproute2",
	"qemu-kvm",
	"qemu-utils",
}

func InstallAgent(ctx *cli.Context) error {
	args := ctx.String("arguments")
	target := ctx.String("base-directory")

	d, err := deploy.NewLocalDeployer(target)
	if err != nil {
		return errors.Wrap(err, "Failed to create new LocalDeployer")
	}

	fmt.Printf("---> [INSTALL] Installing packages: %s\n", strings.Join(AgentRequiredPackages, ", "))
	if err := d.InstallPackages(AgentRequiredPackages, os.Stdout, os.Stderr); err != nil {
		return errors.Wrap(err, "Failed to install packages")
	}

	fmt.Printf("---> [INSTALL] Set sysctl: net.ipv4.ip_forward=1\n")
	if err := d.SetSysctl("net.ipv4.ip_forward", []byte("1")); err != nil {
		return errors.Wrap(err, "Failed to set sysctl")
	}

	fmt.Println("---> [INSTALL] Stopping systemd unit...")
	if err := d.StopDaemon(systemdAgentUnit, os.Stdout, os.Stderr); err != nil {
		return errors.Wrap(err, "Failed to stop systemd daemon")
	}

	binLocation := "/usr/bin/n0core"
	fmt.Printf("---> [INSTALL] Linking self to %s...\n", binLocation)
	if err := d.LinkSelf(binLocation); err != nil {
		return errors.Wrap(err, "Failed to link self")
	}

	fmt.Println("---> [INSTALL] Preparing systemd unit...")
	cmd := fmt.Sprintf("%s serve agent %s", binLocation, args)
	systemd := d.CreateAgentUnit(cmd)
	if err := d.SaveFile(systemd, systemdAgentUnitPath, 0644); err != nil {
		return errors.Wrap(err, "Failed to save systemd unit file")
	}

	fmt.Println("---> [INSTALL] Restarting systemd unit...")
	if err := d.RestartDaemon(systemdAgentUnit, os.Stdout, os.Stderr); err != nil {
		return errors.Wrap(err, "Failed to restart systemd daemon")
	}

	fmt.Println("---> [INSTALL] Waiting 1 secs to show status...")
	time.Sleep(1 * time.Second)
	if err := d.DaemonStatus(systemdAgentUnit, os.Stdout, os.Stderr); err != nil {
		return errors.Wrap(err, "Failed to get status of systemd daemon")
	}

	return nil
}

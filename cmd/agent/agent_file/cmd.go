package main

import (
	loggingGateway "github.com/n0stack/n0core/gateway/logging"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/processor/agent"
	"github.com/n0stack/n0core/target"
	"github.com/n0stack/n0core/target/network/flat"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	// app := cli.NewApp()

	// app.Name = "n0core agent which get message by file"
	// app.Usage = "test test"
	// app.Version = "0.1.0"

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "file, f",
	// 		Value: "test.yml",
	// 		Usage: "",
	// 	},
	// 	cli.StringFlag{
	// 		Name:  "interface, i",
	// 		Value: "eth0",
	// 		Usage: "Interface to connect external network by internal bridge",
	// 	},
	// }

	// buf, err := ioutil.ReadFile("test.yml")
	// if err != nil {
	// 	panic(err)
	// }

	buf := []byte(`taskID: 2efbfd8d-6136-4390-a513-033e7c5f2391
task: Delete
model:
  id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
  type: resource/network/flat
  name: hogehoge
  state: UP
  subnets:
  - cidr: 192.168.0.0/24
    dhcp:
      rangeStart: 192.168.0.1
      rangeEnd: 192.168.0.127
      nameservers:
        - 192.168.0.254
      gateway: 192.168.0.254
`)

	t := message.Task{}
	err := yaml.Unmarshal(buf, &t)
	if err != nil {
		panic(err)
	}

	println("*---- Task message ----*")
	println(string(buf))
	println()

	f := &flat.Flat{InterfaceName: "enp0s25"}
	g := &loggingGateway.LoggingGateway{}

	a, err := agent.NewAgent([]target.Target{f}, g, map[string]string{})
	if err != nil {
		panic(err)
	}

	println("*---- Notification messages ----*")
	err = a.ProcessMessage(&t)
	if err != nil {
		panic(err)
	}
}

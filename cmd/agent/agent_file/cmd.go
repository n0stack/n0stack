package main

import (
	"flag"
	"io/ioutil"

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

	filename := flag.String("f", "", "file name of yaml file which is read as task message.")
	flag.Parse()
	if *filename == "" {
		panic("filename is required.")
	}

	buf, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}

	t := message.Task{}
	err = yaml.Unmarshal(buf, &t)
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

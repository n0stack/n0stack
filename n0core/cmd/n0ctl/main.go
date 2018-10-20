package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0ctl"
	// app.Usage = ""
	// app.Version = "0.1.2" // CIで取るようにする

	app.Commands = []cli.Command{
		{
			Name:      "get",
			Usage:     "Get resource if set resource name, List resources if not set",
			ArgsUsage: "[resource type] [resource name]",
			Action:    Get,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "api-endpoint",
					Value:  "localhost:20180",
					EnvVar: "N0CTL_API_ENDPOINT",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete resource",
			ArgsUsage: "[resource type] [resource name]",
			Action:    Delete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "api-endpoint",
					Value:  "localhost:20180",
					EnvVar: "N0CTL_API_ENDPOINT",
				},
			},
		},
		{
			Name:  "do",
			Usage: "Do DAG tasks (Detail n0stack/pkg/dag)",
			Description: `
	## File format

	---
	task_name:
		resource_type: Network
		action: GetNetwork
		args:
			name: test-network
		depend_on:
			- dependency_task_name
	dependency_task_name:
		resource_type: ...
	---

	- task_name
			- 任意の名前をつけ、ひとつのリクエストに対してユニークなものにする
	- resource_type
			- gRPC メッセージを指定する
			- VirtualMachine や virtual_machine という形で指定できる
	- action
			- gRPC の RPC を指定する
			- GetNetwork など定義のとおりに書く
	- args
			- gRPC の RPCのリクエストを書く
	- depend_on
			- DAG スケジューリングに用いられる
			- task_name を指定する
			`,
			ArgsUsage: "[file name]",
			Action:    Do,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "api-endpoint",
					Value:  "localhost:20180",
					EnvVar: "N0CTL_API_ENDPOINT",
				},
			},
		},
	}

	log.SetFlags(log.Lshortfile)

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to command: %v", err.Error())
	}
}

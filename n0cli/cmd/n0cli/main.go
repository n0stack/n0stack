package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	grpccmd "github.com/n0stack/n0stack/n0cli/grpc_cmd"
)

var version = "undefined"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	outputter := grpccmd.DefaultOutputer()

	app := cli.NewApp()
	app.Name = "n0cli"
	app.Version = version
	app.EnableBashCompletion = true
	app.Usage = "the n0stack CLI application"
	app.Description = `
	## Bash Completion

	---
	wget https://raw.githubusercontent.com/urfave/cli/master/autocomplete/bash_autocomplete
	PROG=n0cli source bash_autocomplete
	---
	`

	app.Flags = []cli.Flag{
		grpccmd.API_URL_FLAG,
		cli.StringFlag{
			Name:   "api-endpoint",
			Value:  "localhost:20180",
			EnvVar: "N0CLI_API_ENDPOINT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "do",
			Usage: "Do DAG tasks (Detail n0stack/pkg/dag)",
			Description: `
	## File format

	---
	task_name:
		type: Network
		action: GetNetwork
		args:
			name: test-network
		depend_on:
			- dependency_task_name
		ignore_error: true
	dependency_task_name:
		type: ...
	---

	- task_name
			- 任意の名前をつけ、ひとつのリクエストに対してユニークなものにする
	- type
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
	- ignore_error
	    - タスクでエラーが発生しても継続する
			`,
			ArgsUsage: "[file name]",
			Action:    Do,
		},
		{
			Name:        "node",
			Usage:       "Node APIs",
			Description: "",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get Node(s)",
					ArgsUsage: "[Node name (optional) ...]",
					Action:    GetNode,
				},
				{
					Name:      "delete",
					Usage:     "Delete Node",
					ArgsUsage: "[Node name]",
					Action:    DeleteNode,
				},
			},
		},
		{
			Name:        "network",
			Aliases:     []string{"net"},
			Usage:       "Network APIs",
			Description: "",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get Network(s)",
					ArgsUsage: "[Network name (optional) ...]",
					Action: func(c *cli.Context) error {
						out := outputter.GenerateOutputMethod([]string{"name", "ipv4_cidr", "ipv6_cidr"})
						if c.NArg() == 0 {
							f := grpccmd.GenerateAction(ctx, out, ppool.NewNetworkServiceClient, ppool.NetworkServiceClient.ListNetworks, []string{})
							return f(c)
						} else if c.NArg() == 1 {
							f := grpccmd.GenerateAction(ctx, out, ppool.NewNetworkServiceClient, ppool.NetworkServiceClient.GetNetwork, []string{"name"})
							return f(c)
						}

						return fmt.Errorf("set valid arguments")
					},
					Flags: append(grpccmd.GenerateFlags(ppool.NetworkServiceClient.GetNetwork), grpccmd.OUTPUT_TYPE_FLAG),
				},
				{
					Name:      "delete",
					Usage:     "Delete Network",
					ArgsUsage: "[Network name]",
					Action:    DeleteNetwork,
				},
				{
					Name:      "apply",
					Usage:     "Apply Network",
					ArgsUsage: "[Network name]",
					Action:    grpccmd.GenerateAction(ctx, outputter.OutputJsonAsOutputMessage, ppool.NewNetworkServiceClient, ppool.NetworkServiceClient.ApplyNetwork, []string{"name"}),
					Flags:     grpccmd.GenerateFlags(ppool.NetworkServiceClient.ApplyNetwork),
				},
			},
		},

		{
			Name:        "virtual_machine",
			Aliases:     []string{"vm"},
			Usage:       "VirtualMachine APIs",
			Description: "",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get VirtualMachine(s)",
					ArgsUsage: "[VirtualMachine name (optional) ...]",
					Action:    GetVirtualMachine,
				},
				{
					Name:      "delete",
					Usage:     "Delete VirtualMachine",
					ArgsUsage: "[VirtualMachine name]",
					Action:    DeleteVirtualMachine,
				},
				{
					Name:      "boot",
					Usage:     "Boot VirtualMachine",
					ArgsUsage: "[VirtualMachine name]",
					Action:    BootVirtualMachine,
				},
				{
					Name:      "open_console",
					Aliases:   []string{"console"},
					Usage:     "Get URL to open console of VirtualMachine",
					ArgsUsage: "[VirtualMachine name]",
					Action:    OpenConsoleOfVirtualMachine,
				},
			},
		},
		{
			Name:        "block_storage",
			Aliases:     []string{"bs"},
			Usage:       "BlockStorage APIs",
			Description: "",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get BlockStorage(s)",
					ArgsUsage: "[BlockStorage name (optional) ...]",
					Action:    GetBlockStorage,
				},
				{
					Name:      "delete",
					Usage:     "Delete BlockStorage",
					ArgsUsage: "[BlockStorage name]",
					Action:    DeleteBlockStorage,
				},
				{
					Name:      "download",
					Usage:     "Get URL to download BlockStorage",
					ArgsUsage: "[BlockStorage name]",
					Action:    DownloadBlockStorage,
				},
			},
		},
		{
			Name:        "image",
			Usage:       "Image APIs",
			Description: "",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get Image(s)",
					ArgsUsage: "[Image name (optional) ...]",
					Action:    GetImage,
				},
				{
					Name:      "delete",
					Usage:     "Delete Image",
					ArgsUsage: "[Image name]",
					Action:    DeleteImage,
				},
				{
					Name:        "register",
					Usage:       "Add [BlockStorage] to [Image name] WITH ApplyImageRequest [Image name]",
					Description: "register blockstorage to image.",
					ArgsUsage:   "-a [Image name] [BlockStorage name]",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "apply-image, a",
							Usage: "If [Image name] doesn't exist, it send ApplyImageRequest.",
						},
					},
					Action: RegisterBlockStorage,
				},
				{
					Name:      "unregister",
					Usage:     "Remove [BlockStorage name] from [Image name]",
					ArgsUsage: "[Image name] [BlockStorage name]",
					Action:    UnregisterBlockStorage,
				},
				{
					Name:      "tag",
					Usage:     "Add [Tag name] to [BlockStorage name] of [Image name]",
					ArgsUsage: "[Image name]:[Tag name] [BlockStorage name]",
					Action:    Tag,
				},
				{
					Name:      "untag",
					Usage:     "Remove [Tag name] from [Image name]",
					ArgsUsage: "[Image name]:[Tag name]",
					Action:    Untag,
				},
			},
		},
	}

	getCommand := cli.Command{
		Name:      "get",
		Usage:     "Get resource(s)",
		ArgsUsage: "[resource name (optional) ...]",
	}
	deleteCommand := cli.Command{
		Name:      "delete",
		Usage:     "Delete resource",
		ArgsUsage: "[resource name]",
	}
	for _, c1 := range app.Commands {
		if c1.Name == "do" {
			continue
		}

		for _, c2 := range c1.Subcommands {
			if c2.Name == "get" {
				getCommand.Subcommands = append(getCommand.Subcommands, cli.Command{
					Name:      c1.Name,
					Aliases:   c1.Aliases,
					Usage:     c2.Usage,
					ArgsUsage: c2.ArgsUsage,
					Action:    c2.Action,
					Flags:     c2.Flags,
				})
			} else if c2.Name == "delete" {
				deleteCommand.Subcommands = append(deleteCommand.Subcommands, cli.Command{
					Name:      c1.Name,
					Aliases:   c1.Aliases,
					Usage:     c2.Usage,
					ArgsUsage: c2.ArgsUsage,
					Action:    c2.Action,
					Flags:     c2.Flags,
				})
			}
		}
	}
	app.Commands = append(app.Commands, getCommand)
	app.Commands = append(app.Commands, deleteCommand)

	log.SetFlags(log.Lshortfile)
	log.SetOutput(ioutil.Discard)

	if err := app.Run(os.Args); err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "Failed to command: %s\n", err.Error())
		color.Unset()
		os.Exit(1)
	}
}

func ConnectAPI(c *cli.Context) (*grpc.ClientConn, error) {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	return conn, nil
}

func PrintGrpcError(err error) {
	fmt.Fprintf(os.Stderr, "[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
}

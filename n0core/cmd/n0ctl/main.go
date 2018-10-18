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
			Usage:     "Get resource if set resource_name, List resources if not set",
			ArgsUsage: "[resource_type] [resource_name]",
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
			ArgsUsage: "[resource_type] [resource_name]",
			Action:    Delete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "api-endpoint",
					Value:  "localhost:20180",
					EnvVar: "N0CTL_API_ENDPOINT",
				},
			},
		},
		// {
		// 	Name:   "do",
		// 	Usage:  "",
		// 	Action: Do,
		// },
	}

	log.SetFlags(log.Lshortfile)

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to command: err=%v", err.Error())
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/n0stack/n0stack/n0proto.go/pkg/dag"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

// TODO: 複数ファイルを連結して処理できるようにしたい `n0ctl do foo.yaml bar.yaml` みたいな
func Do(ctx *cli.Context) error {
	if ctx.NArg() == 1 {
		return do(ctx)
	}

	return fmt.Errorf("set valid arguments")
}

// TODO: エラーレスポンス
func do(ctx *cli.Context) error {
	filepath := ctx.Args().Get(0)
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	tasks := map[string]*dag.Task{}
	// _, err := toml.DecodeFile(filepath, &tasks)
	if err := yaml.Unmarshal(buf, tasks); err != nil {
		return err
	}

	endpoint := ctx.String("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	dag.Marshaler = marshaler
	if err := dag.CheckDAG(tasks); err != nil {
		return err
	}

	if ok := dag.DoDAG(tasks, os.Stdout, conn); !ok {
		return fmt.Errorf("Failed to do tasks")
	}

	fmt.Fprintf(os.Stderr, "DAG tasks are completed\n")
	return nil
}

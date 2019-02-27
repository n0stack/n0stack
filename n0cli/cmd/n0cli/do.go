package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/n0stack/n0stack/n0proto.go/pkg/dag"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

// TODO: 複数ファイルを連結して処理できるようにしたい `n0cli do foo.yaml bar.yaml` みたいな
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
	if err := yaml.UnmarshalStrict(buf, tasks); err != nil {
		return err
	}

	endpoint := ctx.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	ctxCancel, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
	}()

	go func() {
		select {
		case <-c: // SIGINT
			cancel()       // notify DoDAG to cancel
			signal.Stop(c) // allow sending SIGINT again to force SIGINT
		case <-ctxCancel.Done():
			return
		}
	}()

	dag.Marshaler = marshaler
	if err := dag.CheckDAG(tasks); err != nil {
		return err
	}

	if ok := dag.DoDAG(ctxCancel, tasks, os.Stdout, conn); !ok {
		return fmt.Errorf("Failed to do tasks")
	}

	fmt.Fprintf(os.Stderr, "DAG tasks are completed\n")
	return nil
}

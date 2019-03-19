package grpccmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var marshaler = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "  ",
	OrigName:     true,
}

var OUTPUT_TYPE_FLAG = cli.StringFlag{
	Name:  "output, o",
	Value: "table",
	Usage: "Select a output method from `table,json`",
}

type Outputter struct {
	Out io.Writer
	Err io.Writer
}

type OutputMessage func(*cli.Context, proto.Message) error

func DefaultOutputer() *Outputter {
	return &Outputter{
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

func (o Outputter) GenerateOutputMethod(tableKeys []string) OutputMessage {
	return func(c *cli.Context, m proto.Message) error {
		output := c.String("output")

		switch output {
		case "table":
			return o.OutputTable(m, tableKeys)

		case "json":
			return o.OutputJson(m)

			// case "yaml":
			// 	return o.OutputYaml(m)
		}

		return errors.New("invalid output type")
	}
}

func (o Outputter) OutputNone(c *cli.Context, m proto.Message) error {
	return nil
}

func (o Outputter) OutputJsonAsOutputMessage(c *cli.Context, m proto.Message) error {
	return o.OutputJson(m)
}

func (o Outputter) OutputJson(m proto.Message) error {
	if err := marshaler.Marshal(o.Out, m); err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func (o Outputter) OutputTable(m proto.Message, keys []string) error {
	buf := &bytes.Buffer{}
	err := marshaler.Marshal(buf, m)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal message to json")
	}

	value := make(map[string]interface{})
	if err := json.Unmarshal(buf.Bytes(), &value); err != nil {
		return errors.Wrap(err, "")
	}

	var data [][]string
	if _, ok := value["name"]; ok { // all resources have "name" key. if value do not have "name", it is a list response.
		data = make([][]string, 1)
		data[0] = mapToData(value, keys)
	} else {
		for _, v := range value {
			buf, err := json.Marshal(v)
			if err != nil {
				return errors.Wrap(err, "")
			}

			v := make([]map[string]interface{}, 0)
			if err := json.Unmarshal(buf, &v); err != nil {
				return errors.Wrap(err, "")
			}
			data = make([][]string, len(v))

			for i, value := range v {
				data[i] = mapToData(value, keys)
			}

			break // first value is listed resources
		}
	}

	table := tablewriter.NewWriter(o.Out)
	table.SetHeader(keys)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()

	fmt.Println()
	return nil
}

func mapToData(value map[string]interface{}, keys []string) []string {
	d := make([]string, len(keys))
	for j, k := range keys {
		if v, ok := value[k]; ok {
			b, _ := json.Marshal(v)
			d[j] = strings.TrimLeft(strings.TrimRight(string(b), "\""), "\"")
		} else {
			d[j] = ""
		}
	}

	return d
}

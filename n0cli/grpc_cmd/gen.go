package grpccmd

import (
	"context"
	"reflect"
	"strings"

	"github.com/urfave/cli"
)

func ParseJsonTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}

	return tag
}

// "name" is standard field, so get by args
func GenerateFlags(targetGRPC interface{}) []cli.Flag {
	t := reflect.TypeOf(targetGRPC).In(2).Elem()
	flags := []cli.Flag{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := ParseJsonTag(field.Tag.Get("json"))
		if tag == "-" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			flags = append(flags, cli.StringFlag{Name: tag})

		case reflect.Int64:
			flags = append(flags, cli.Int64Flag{Name: tag})
		case reflect.Int32:
			flags = append(flags, cli.IntFlag{Name: tag})
		case reflect.Uint64:
			flags = append(flags, cli.Uint64Flag{Name: tag})
		case reflect.Uint32:
			flags = append(flags, cli.UintFlag{Name: tag})

		case reflect.Slice:
			// []string, []structに対応
			flags = append(flags, cli.StringSliceFlag{Name: tag})

			// TODO: []int など
			// if field.Type == reflect.TypeOf([]string{}) {
			// }

		case reflect.Map:
			// --map=key:value のような使い方を想定
			flags = append(flags, cli.StringSliceFlag{
				Name:  tag,
				Usage: "set like --option=[key]:[value]",
			})
		}
	}

	return flags
}

func GenerateAction(ctx context.Context, output OutputMessage, newGrpcClient interface{}, targetGRPC interface{}, argsKeys []string) func(*cli.Context) error {
	getter := GenerateGRPCGetter(targetGRPC, argsKeys, newGrpcClient)

	return func(c *cli.Context) error {
		conn, err := Connect2gRPC(c)
		if err != nil {
			return err
		}
		defer conn.Close()

		m, err := getter(c, ctx, conn)
		if err != nil {
			return err
		}

		if err := output(c, m); err != nil {
			return err
		}

		return nil
	}
}

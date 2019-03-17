package grpccmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

var API_URL_FLAG = cli.StringFlag{
	Name:   "api-url",
	Value:  "grpc://localhost:20180",
	EnvVar: "N0CLI_API_URL",
}

func Connect2gRPC(c *cli.Context) (*grpc.ClientConn, error) {
	endpoint := c.GlobalString("api-url")
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(u.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	return conn, nil
}

// this function panic when set an argument that is not gRPC method.
func GenerateGRPCGetter(f interface{}, argsKeys []string) func(c *cli.Context, ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	reqType := t.In(2) // reqType means *RequestMessage

	return func(c *cli.Context, ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
		if c.NArg() != len(argsKeys) {
			return nil, fmt.Errorf("set valid arguments")
		}

		for i, a := range argsKeys {
			c.Set(a, c.Args()[i])
			log.Println(c.String(a))
		}

		req := reflect.New(reqType.Elem())           // New() reqturns generated argument pointer, so New(*RequestMessage.Elem()) -> *RequestMessage{}
		for i := 0; i < req.Elem().NumField(); i++ { // Call NumFiled on RequestMessage{}
			field := reqType.Elem().Field(i)
			tag := ParseJsonTag(field.Tag.Get("json"))
			if tag == "-" {
				continue
			}

			target := req.Elem().Field(i) // 多分大丈夫だが少し不安
			switch field.Type.Kind() {
			case reflect.String:
				target.Set(reflect.ValueOf(c.String(tag)))

			case reflect.Int64:
				target.Set(reflect.ValueOf(c.Int64(tag)))
			case reflect.Int32:
				target.Set(reflect.ValueOf(c.Int(tag)))
			case reflect.Uint64:
				target.Set(reflect.ValueOf(c.Uint64(tag)))
			case reflect.Uint32:
				target.Set(reflect.ValueOf(c.Uint(tag)))

			case reflect.Slice:
				target.Set(reflect.ValueOf(c.StringSlice(tag)))

			case reflect.Map:
				s := c.StringSlice(tag)
				m := make(map[string]string)

				for _, t := range s {
					kv := strings.Split(t, ":")
					if len(kv) > 2 {
						m[kv[0]] = strings.Join(kv[1:], ":")
					} else if len(kv) == 2 {
						m[kv[0]] = kv[1]
					} else {
						return nil, fmt.Errorf("Failed to parse map argument: --%s=%s", tag, t)
					}
				}
				target.Set(reflect.ValueOf(m))
			}
		}

		log.Printf("[DEBUG] request: %+v", req)
		cli := ppool.NewNetworkServiceClient(conn) // TODO: ここをうまくやる
		out := v.Call([]reflect.Value{reflect.ValueOf(cli), reflect.ValueOf(ctx), req})
		if err, _ := out[1].Interface().(error); err != nil {
			return nil, fmt.Errorf("got error response: %s", err.Error())
		}

		return out[0].Interface().(proto.Message), nil
	}
}

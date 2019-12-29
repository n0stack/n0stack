package grpccmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"n0st.ac/n0stack/n0core/pkg/driver/n0stack/auth"
	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"
	structutil "n0st.ac/n0stack/n0core/pkg/util/struct"
)

var API_URL_FLAG = cli.StringFlag{
	Name:   "api-url",
	Value:  "grpc://localhost:20180",
	EnvVar: "N0CLI_API_URL",
}
var API_INSECURE_FLAG = cli.BoolFlag{
	Name: "insecure,k",
}
var IDENTITY_FILE_FLAG = cli.StringFlag{
	Name: "identity-file,i",
}
var LOGIN_NAME_FLAG = cli.StringFlag{
	Name: "login-name,l",
}

func Connect2gRPC(c *cli.Context) (*grpc.ClientConn, error) {
	endpoint := c.GlobalString("api-url")
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	host := u.Host
	if u.Port() == "" {
		host += ":443"
	}
	if strings.HasSuffix(host, ":443") {
		if !c.GlobalBool("insecure") {
			pool, err := x509.SystemCertPool()
			if err != nil {
				log.Fatal("failed to get system certifications")
			}

			creds := credentials.NewClientTLSFromCert(pool, "")
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	identityFile := c.GlobalString("identity-file")
	if identityFile != "" {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "grpc.Dial(%s, %+v) returns err=%+v to get a token", host, opts, err)
		}
		defer conn.Close()

		user := c.GlobalString("login-name")
		key, err := jwtutil.ParsePrivateKeyFromFile(identityFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse the private key %s", identityFile)
		}

		authc, err := auth.NewAuthenticationClient(context.Background(), conn, user, key)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare a authentication token")
		}

		opts = append(opts, grpc.WithPerRPCCredentials(authc))
	}

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	return conn, nil
}

// this function panic when set an argument that is not gRPC method.
func GenerateGRPCGetter(f interface{}, argsKeys []string, newGrpcClient interface{}) func(c *cli.Context, ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	reqType := t.In(2) // reqType means *RequestMessage

	return func(c *cli.Context, ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
		if c.NArg() != len(argsKeys) {
			return nil, fmt.Errorf("set valid arguments, got wrong number of arguments: got=%d, have=%d", c.NArg(), len(argsKeys))
		}

		for i, a := range argsKeys {
			c.Set(a, c.Args()[i])
			log.Println(c.String(a))
		}

		req := reflect.New(reqType.Elem()) // New() reqturns generated argument pointer, so New(*RequestMessage.Elem()) -> *RequestMessage{}

		for _, f := range c.FlagNames() {
			if f == strings.Split(OUTPUT_TYPE_FLAG.Name, ",")[0] {
				continue
			}

			v, err := structutil.GetValueByJson(req, f)
			if err != nil {
				return nil, errors.Wrapf(err, "GetValueByJson(%+v, %s) returns err=%+v", req, f, err)
			}

			var value interface{}
			switch v.Kind() {
			case reflect.String:
				value = c.String(f)

			case reflect.Bool:
				value = c.Bool(f)

			case reflect.Int64:
				value = c.Int64(f)
			case reflect.Int32:
				value = c.Int(f)
			case reflect.Uint64:
				value = c.Uint64(f)
			case reflect.Uint32:
				value = c.Uint(f)

			case reflect.Slice:
				value = c.StringSlice(f)

			case reflect.Map:
				s := c.StringSlice(f)
				m := make(map[string]string)

				for _, t := range s {
					kv := strings.Split(t, ":")
					if len(kv) > 2 {
						m[kv[0]] = strings.Join(kv[1:], ":")
					} else if len(kv) == 2 {
						m[kv[0]] = kv[1]
					} else {
						return nil, fmt.Errorf("Failed to parse map argument: --%s=%s", f, t)
					}
				}

				value = m
			}

			if value == nil {
				continue
			}

			if err := structutil.SetByJson(req.Interface(), f, value); err != nil {
				return nil, errors.Wrapf(err, "failed to set %+v as request filed to %s", value, f)
			}
		}

		log.Printf("[DEBUG] request: %+v", req.Interface())
		newCli := reflect.ValueOf(newGrpcClient)
		cli := newCli.Call([]reflect.Value{reflect.ValueOf(conn)})[0]
		out := v.Call([]reflect.Value{cli, reflect.ValueOf(ctx), req})
		if err, _ := out[1].Interface().(error); err != nil {
			PrintGrpcError(err)
			return nil, fmt.Errorf("")
		}

		return out[0].Interface().(proto.Message), nil
	}
}

func PrintGrpcError(err error) {
	fmt.Fprintf(os.Stderr, "[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
}

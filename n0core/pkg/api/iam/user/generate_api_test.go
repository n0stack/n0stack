package userapi

import (
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/util/generator"
)

// const stringOneRun = `func (i %[2]s)

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const ListUsersTemplate = `func List%[2]ss(ctx context.Context, req *p%[1]s.ListUsersRequest, ds datastore.Datastore) (*p%[1]s.List%[2]ssResponse, error) {
	res := &p%[1]s.List%[2]ssResponse{}
	f := func(s int) []proto.Message {
		res.%[2]ss = make([]*p%[1]s.%[2]s, s)
		for i := range res.%[2]ss {
			res.%[2]ss[i] = &p%[1]s.%[2]s{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.%[2]ss {
			m[i] = v
		}

		return m
	}

	if err := ds.List(f); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.%[2]ss) == 0 {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	return res, nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const GetUserTemplate = `func Get%[2]s(ctx context.Context, req *p%[1]s.Get%[2]sRequest, ds datastore.Datastore) (*p%[1]s.%[2]s, error) {
	resourse := &p%[1]s.%[2]s{}
	if err := ds.Get(req.Name, resourse); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	return resourse, nil
}`

func GenerateListAPI(gen *generator.GoCodeGenerator, service, resource string) {
	gen.AddPackage("", "context")
	gen.AddPackage("", "github.com/golang/protobuf/proto")
	gen.AddPackage("", "github.com/n0stack/n0stack/n0core/pkg/datastore")
	gen.AddPackage("grpcutil", "github.com/n0stack/n0stack/n0core/pkg/util/grpc")
	gen.AddPackage("piam", "github.com/n0stack/n0stack/n0proto.go/iam/v0")
	gen.AddPackage("", "google.golang.org/grpc/codes")

	gen.Printf(ListUsersTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")
}

func GenerateGetAPI(gen *generator.GoCodeGenerator, service, resource string) {
	gen.AddPackage("", "context")
	gen.AddPackage("", "github.com/golang/protobuf/proto")
	gen.AddPackage("", "github.com/n0stack/n0stack/n0core/pkg/datastore")
	gen.AddPackage("grpcutil", "github.com/n0stack/n0stack/n0core/pkg/util/grpc")
	gen.AddPackage("piam", "github.com/n0stack/n0stack/n0proto.go/iam/v0")
	gen.AddPackage("", "google.golang.org/grpc/codes")

	gen.Printf(GetUserTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")
}

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("template_api", "userapi")
	GenerateListAPI(g, "iam", "User")
	GenerateGetAPI(g, "iam", "User")

	if err := g.WriteWithTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}

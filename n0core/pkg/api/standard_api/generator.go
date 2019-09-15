package stdapi

import "github.com/n0stack/n0stack/n0core/pkg/util/generator"

// Arguments to format are:
//	[1]: package
//	[2]: package version
const ImportedTemplate = `import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	p%[1]s "github.com/n0stack/n0stack/n0proto.go/%[1]s/%[2]s"
	"google.golang.org/grpc/codes"
)`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const ListTemplate = `func List%[2]ss(ctx context.Context, req *p%[1]s.List%[2]ssRequest, ds datastore.Datastore) (*p%[1]s.List%[2]ssResponse, error) {
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

	if err := ds.List(ctx, f); err != nil {
		return nil, grpcutil.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.%[2]ss) == 0 {
		return nil, grpcutil.Errorf(codes.NotFound, "")
	}

	return res, nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const GetTemplate = `func Get%[2]s(ctx context.Context, ds datastore.Datastore, name string) (*p%[1]s.%[2]s, int64, error) {
	resourse := &p%[1]s.%[2]s{}
	version, err := ds.Get(ctx, name, resourse)
	if err != nil {
		if datastore.IsNotFound(err) {
			return nil, 0, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, 0, grpcutil.Errorf(codes.Internal, "failed to get User %%s from db: err='%%s'", name, err.Error())
	}

	return resourse, version, nil
}`

// Arguments to format are:
//	[1]: resource name
const DeleteTemplate = `func Delete%[1]s(ctx context.Context, ds datastore.Datastore, name string, version int64) error {
	if err := ds.Delete(ctx, name, version); err != nil {
		return grpcutil.Errorf(codes.Internal, "failed to delete %[1]s %%s from db: err='%%s'", name, err.Error())
	}

	return nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const ApplyTemplate = `func Apply%[2]s(ctx context.Context, ds datastore.Datastore, resource *p%[1]s.%[2]s, version int64) (int64, error) {
	version, err := ds.Apply(ctx, resource.Name, resource, version);
	if err != nil {
		return 0, grpcutil.Errorf(codes.Internal, "failed to apply %[2]s %%s to db: err='%%s'", resource.Name, err.Error())
	}

	return version, nil
}`

func GenerateTemplateAPI(gen *generator.GoCodeGenerator, service, resource, version string) {
	gen.Printf(ImportedTemplate, service, version)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(GetTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(DeleteTemplate, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(ApplyTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")
}

func GenerateTemplatedListAPI(gen *generator.GoCodeGenerator, service, resource string) {
	gen.Printf(ListTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")
}

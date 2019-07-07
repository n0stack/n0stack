package stdapi

import "github.com/n0stack/n0stack/n0core/pkg/util/generator"

// Arguments to format are:
//	[1]: package
const ImportedTemplate = `import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	p%[1]s "github.com/n0stack/n0stack/n0proto.go/%[1]s/v0"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
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
const GetTemplate = `func Get%[2]s(ctx context.Context, req *p%[1]s.Get%[2]sRequest, ds datastore.Datastore) (*p%[1]s.%[2]s, error) {
	resourse := &p%[1]s.%[2]s{}
	if err := ds.Get(req.Name, resourse); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	return resourse, nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const GetAndPendExistingTemplate = `func GetAndPendExisting%[2]s(tx *transaction.Transaction, ds datastore.Datastore, name string) (*p%[1]s.%[2]s, error) {
	resource := &p%[1]s.%[2]s{}
	if err := ds.Get(name, resource); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	if resource.State == p%[1]s.%[2]s_PENDING {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "%[2]s %%s is pending", name)
	}

	current := resource.State
	resource.State = p%[1]s.%[2]s_PENDING
	if err := Apply%[2]s(ds, resource); err != nil {
		return nil, err
	}
	resource.State = current
	tx.PushRollback("free optimistic lock", func() error {
		resource.State = current
		return ds.Apply(resource.Name, resource)
	})

	return resource, nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const PendNewTemplate = `func PendNew%[2]s(tx *transaction.Transaction, ds datastore.Datastore, name string) error {
	resource := &p%[1]s.%[2]s{}
	if err := ds.Get(name, resource); err == nil {
		return grpcutil.WrapGrpcErrorf(codes.AlreadyExists, "%[2]s %%s is already exists", name)
	} else if !datastore.IsNotFound(err) {
		return grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	resource.Name = name
	resource.State = p%[1]s.%[2]s_PENDING
	if err := Apply%[2]s(ds, resource); err != nil {
		return err
	}
	tx.PushRollback("free optimistic lock", func() error {
		return ds.Delete(name)
	})

	return nil
}`

// Arguments to format are:
//	[1]: resource name
const DeleteTemplate = `func Delete%[1]s(ds datastore.Datastore, name string) error {
	if err := ds.Delete(name); err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "failed to delete %[1]s %%s from db: err='%%s'", name, err.Error())
	}

	return nil
}`

// Arguments to format are:
//	[1]: package
//	[2]: resource name
const ApplyTemplate = `func Apply%[2]s(ds datastore.Datastore, resource *p%[1]s.%[2]s) error {
	if err := ds.Apply(resource.Name, resource); err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "failed to apply %[2]s %%s to db: err='%%s'", resource.Name, err.Error())
	}

	return nil
}`

func GenerateTemplateAPI(gen *generator.GoCodeGenerator, service, resource string) {
	gen.Printf(ImportedTemplate, service)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(ListTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(GetTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(GetAndPendExistingTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(PendNewTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(DeleteTemplate, resource)
	gen.Printf("\n")
	gen.Printf("\n")

	gen.Printf(ApplyTemplate, service, resource)
	gen.Printf("\n")
	gen.Printf("\n")
}

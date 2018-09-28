package provisioning

import (
	"context"

	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
)

func (a BlockStorageAPI) reserveStorage(name string, annotations map[string]string, req, limit uint64) (string, string, error) {
	var rcr *ppool.ReserveStorageResponse
	var err error
	if node, ok := annotations[AnnotationRequestNodeName]; !ok {
		rcr, err = a.nodeAPI.ScheduleStorage(context.Background(), &ppool.ScheduleStorageRequest{
			StorageName: name,
			Storage: &pbudget.Storage{
				RequestBytes: req,
				LimitBytes:   limit,
			},
		})
	} else {
		rcr, err = a.nodeAPI.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
			Name:        node,
			StorageName: name,
			Storage: &pbudget.Storage{
				RequestBytes: req,
				LimitBytes:   limit,
			},
		})
	}
	if err != nil {
		return "", "", err // TODO: #89
	}

	return rcr.Name, rcr.StorageName, nil
}

package ablockstorage

import (
	"context"
	"log"

	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	pinfrastructure "n0st.ac/n0stack/n0proto.go/infrastructure/v0"
	"n0st.ac/n0stack/n0proto.go/pkg/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Reserve a block storage on a node using MachineServiceClient.
func ReserveStorage(ctx context.Context, tx *transaction.Transaction, na pinfrastructure.MachineServiceClient, bs *pinfrastructure.BlockStorage) error {
	bs.StorageName = bs.Name

	var n *pinfrastructure.Machine
	var err error
	if node, ok := bs.Annotations[AnnotationBlockStorageRequestMachineName]; !ok {
		n, err = na.ScheduleStorage(ctx, &pinfrastructure.ScheduleStorageRequest{
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			Bytes: bs.Bytes,
		})
	} else {
		n, err = na.ReserveStorage(ctx, &pinfrastructure.ReserveStorageRequest{
			MachineName: node,
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			Bytes: bs.Bytes,
		})
	}
	if err != nil {
		return grpcutil.Errorf(grpc.Code(err), "Failed to ReserveStorage: desc=%s", grpc.ErrorDesc(err))
	}

	bs.MachineName = n.Name

	tx.PushRollback("release stroage", func() error {
		_, err := na.ReleaseStorage(ctx, &pinfrastructure.ReleaseStorageRequest{
			MachineName: bs.MachineName,
			StorageName: bs.StorageName,
		})

		return err
	})

	return nil
}

// Release a block storage on a node using MachineServiceClient.
func ReleaseStorage(ctx context.Context, tx *transaction.Transaction, na pinfrastructure.MachineServiceClient, bs *pinfrastructure.BlockStorage) error {
	_, err := na.ReleaseStorage(context.Background(), &pinfrastructure.ReleaseStorageRequest{
		MachineName: bs.MachineName,
		StorageName: bs.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", bs.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return grpcutil.Errorf(codes.Internal, "Failed to release compute '%s': please retry", bs.StorageName)
		}
	}
	tx.PushRollback("reserve storage", func() error {
		_, err = na.ReserveStorage(context.Background(), &pinfrastructure.ReserveStorageRequest{
			MachineName: bs.MachineName,
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			Bytes: bs.Bytes,
		})

		return err
	})

	return nil
}

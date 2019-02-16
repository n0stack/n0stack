# Isolation writing

|||
|--|--|
| Status | accepted |

## Context

- 同じオブジェクトに対する更新系のエンドポイントが同時に実行された場合、あとに終了する操作のみが反映されてしまうためDBと実際の状況が不整合を起こしてしまう

## Decision

- 以下の操作を行うことで楽観的なロックを行う
    - リクエストを受けたとき、DB に保存されているオブジェクトを取得し、 PENDING 状態ではないことを確認
    - PENDING 状態のオブジェクトを DB 保存

### Example in BlockStorage

- 作成の場合
    - `PENDING` で保存
    - 失敗した場合は削除

```go
func (a *BlockStorageAPI) CheckAndLock(tx *transaction.Transaction, bs *pprovisioning.BlockStorage) error {
	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(bs.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", bs.Name)
	} else if prev.Name != "" {
		return grpcutil.WrapGrpcErrorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", bs.Name)
	}

	bs.State = pprovisioning.BlockStorage_PENDING
	if err := a.dataStore.Apply(bs.Name, bs); err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}
	tx.PushRollback("free optimistic lock", func() error {
		return a.dataStore.Delete(bs.Name)
	})

	return nil
}
```

- 更新の場合
    - `PENDING` になってないか確認
    - `PENDING` で保存
    - 失敗した場合は前のステートに変更

```go
func (a *BlockStorageAPI) GetAndLock(tx *transaction.Transaction, name string) (*pprovisioning.BlockStorage, error) {
	bs := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(name, bs); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", name)
	} else if bs.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if bs.State == pprovisioning.BlockStorage_PENDING {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "BlockStorage '%s' is pending", name)
	}

	current := bs.State
	bs.State = pprovisioning.BlockStorage_PENDING
	if err := a.dataStore.Apply(bs.Name, bs); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}
	bs.State = current
	tx.PushRollback("free optimistic lock", func() error {
		return a.dataStore.Apply(bs.Name, bs)
	})

	return bs, nil
}
```

## Consequences

| target | updated |
|--|--|
| n0core/pkg/api/pool/node | No |
| n0core/pkg/api/pool/network | Yes |
| n0core/pkg/api/provisioning/block_storage | Yes |
| n0core/pkg/api/provisioning/virtual_machine | Yes |
| n0core/pkg/api/provisioning/Image | No |
| n0core/pkg/api/provisioning/Flavor | No |

- すでに開発済みである以上のもの以外は、開発段階で従うこと

# Pending State

|||
|--|--|
| Status | accepted |

## Context

APIが障害になった場合、どこまで処理を行ったかわからず、不整合の原因になると考えられる。特に、VirtualMachineやBlockStorageなど実体の操作が伴うものはレスポンスまでの時間が長いため、API障害の影響を受けやすいと考えられる。

## Decision

- VirtualMachineやBlockStorageのCreateなど実体の操作が伴うものは最初に `PENDING` ステートに設定
    - `PENDING` ステートのものは更新を行えないようにする

これによって、APIの故障によって処理が止まってものは `PENDING` ステートによって操作がロックされ、不整合の拡大を抑制することができる。管理者は手動で不整合が起きていないか確認を行い、復旧することで正常性を維持する。

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

- `PENDING` のものが多くなってくると運用に耐えられなくなると考えられるので、実際に動かしながら確認
- networkにも組み込んでしまったが本来はいらない
    - 本来は[これ](lock)の目的で実装していたが、全く効果がなかったので理由が変更になった

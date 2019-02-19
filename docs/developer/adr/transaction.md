# Transaction for update process

|||
|--|--|
| Status | accepted |

## Context

- n0stackはgRPCを使って粗結合に実装しているため、マイクロサービスと同様の問題を抱えている
    - 特に途中で処理が失敗したときにもとに戻す必要がある (原子性)

## Decision

- `github.com/n0stack/n0stack/n0proto.go/pkg/transaction` に実装した
- `Transaction` に行った操作の逆の関数をpush指定いき、失敗したときにその関数をシーケンシャルに逆に呼んでいくことでロールバックを実現する

## Consequences

- 適宜更新
- 正直 `github.com/n0stack/n0stack/n0proto.go/pkg/transaction` の場所は失敗したと思っている
- このトランザクションログを他のAPIと共有できれば、障害時にも強くなりそうだが現状思いつかない

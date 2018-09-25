# Node

## Principle

- 行いたいことは以下のとおりである
  - bare metal provisioning (本来の目的であるが実装予定)
  - サーバインベントリ管理
  - n0coreのサービスディスカバリ
  - 死活監視
- nodeやユーザーがサーバインベントリを登録する
- memberlistでサービスの死活監視を行う
- Nodeとmemberlistの両方があることでサービスとして利用できると考える

## Discovery / Alive monitoring

- With [memberlist](https://github.com/hashicorp/memberlist).
- データの優先度は `memberlist > Node model`

### Agentが正常開始 (Node作成 -> memberlist参加)

1. AgentがAPIにNodeを保存、APIはmemberlistにないのでNotReadyに
2. AgentがAPIを起点にmemberlistに参加、APIは通知からNodeをReadyに

### Agentが正常開始 (memberlist参加 -> Node作成)

1. AgentがAPIを起点にmemberlistに参加、APIはNodeがないのでなにも変更せず
2. AgentがAPIにNodeを保存、APIのmemberlistにあるのでReadyに

### Agentが正常終了

1. AgentがAPIからNodeを削除
2. Agentがmemberlistから抜ける

### Agentが異常終了

1. Agentがmemberlistから抜ける
2. APIが離脱を検知し、NodeをNotReadyに

### APIが異常終了 (同時にAgentが異常終了した場合も同様)

- APIが死亡時は動作としては問題がない

#### APIが復活

- TODO: AgentとAPIどちらがどちらのmemberlistにジョインするか
  - Agentからな気がする

### TODO: memberlistとNodeの値が一致しない

- APIへのリクエスト失敗が考えられる
- `NotReady` or 情報の確実性は `memberlist > Node` なのでNodeを更新 or `Invalid`

## Example

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/ListNodes ''
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/GetNode \
'name: "test"'
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/DeleteNode \
'name: "test"'
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/ReserveCompute '
name: "test"
compute_name: "test-reserve"
compute: {
  request_cpu_milli_core: 10
  limit_cpu_milli_core: 10
  request_memory_bytes: 10
  limit_memory_bytes: 10
}'
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/ReleaseCompute '
name: "test"
compute_name: "test-reserve"
'
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/ReserveStorage '
name: "test"
storage_name: "test-reserve"
storage: {
  request_bytes: 10
  limit_bytes: 10
}'
```

```
grpc_cli call localhost:20181 n0stack.pool.NodeService/ReleaseStorage '
name: "test" 
storage_name: "test-reserve" 
'
```


# Data store

## RDBMS

### Version table

- idはn0coreクラスタにおけるリソース全体のバージョンと同義である
- ユーザーのリクエストごとに複数のオブジェクトを１つのレコードで管理する

```yaml
id: 1
created_at: 0000
spec:| # blob yamlとかjsonそのまま
```

### Events table

- 1つのレコードで1つのオブジェクトのイベントを意味する

```yaml
id: 1
version: 1
created_at: 0010
object: 3bbd2bfd-d6ed-4b6e-b863-7085a000d5cc
event: APPLIED
level: SUCCESS
msg: succeeded creating ...
```

#### Events

- SCHEDULED: リソースがどのagentで管理されるかが決定した時に発行する
- APPLIED: リソースがagentによって展開された時に発行する

#### level

- SUCCESS
- FAILURE

### Relations table

### Object table

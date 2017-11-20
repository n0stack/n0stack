# Data store

## RDBMS

### Spec table

- idはn0coreクラスタにおけるリソース全体のバージョンと同義である
- ユーザーのリクエストごとに複数のオブジェクトを１つのレコードで管理する

```yaml
id: integer  # primary
created_at: timestamp
spec: blob  # yamlとかjsonそのまま
```

### Events table

- 1つのレコードで1つのオブジェクトのイベントを意味する

```yaml
id: integer # primary
spec: Spec
created_at: timestamp
object: Object
event: enum{SCHEDULED, APPLIED}
succeeded: bool
msg: string
```

#### Events

- SCHEDULED: リソースがどのagentで管理されるかが決定した時に発行する
- APPLIED: リソースがagentによって展開された時に発行する

### Relations table

```yaml
id: integer  # primary
directed: bool
from: Object
to: Object  # もしくはuuid? 指定先まで時系列を保存しておく場合はObject
property: Maps <string, string>  # json?
```

### Object table

- 各オブジェクトの共通部分
- protobufにかかれているdataと同じ構成になる

```yaml
id: integer  # primary
object_id: UUID  # not_null
name: string
meta: Maps <string, string>  # json?
```

#### vm

```yaml
state: enum
arch: string
vcpus: int
memory: string
vnc_password: string(password)
```

# Objects

オブジェクトは以下の情報を持っている。

```yaml
id: UUID
type: string
state: enum
meta: Maps <string, string>
dependencies:
  - object: Object
    property: Maps <string, string>
```

### meta

- 外部のコンポーネントとの連携など様々な使い方ができると考えている

#### Example 1: n0gateway

以下のようなリクエストをユーザーが行うと `9c2476ab-dc1e-4904-b8a4-6d991fdc7770` のUUIDに関連付けられているLBサービスにportが参加する。

```yaml
type: resuorce/port
state: running
meta:
  n0stack/n0gateway/join: 9c2476ab-dc1e-4904-b8a4-6d991fdc7770
```

n0gatewayとしては `/api/spec` を監視していればサービスディスカバリを用意に実装することができる。

## Data

- [Resources](resources.md)
- [Agents](agents.md)

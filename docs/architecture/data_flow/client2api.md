# Client => API (HTTP)

## 要件

- ユーザーの使いやすいものにする
- 1つのリクエストでいっぺんにデプロイできるようにする
    - vagrant, terraformっぽい?
    - VMに必要なボリュームなどの依存関係を解決する
- すでに作成されているリソースをIDからアタッチできる
- すべての操作をシークして取得できるようにすることで外部コンポーネントを開発しやすくする
- ユーザーはリクエストを送った時には次の操作ができるようにする必要があるので、レスポンスにはIDがふられている必要がある

## 原則・モデル

- `POST /api/spec` でいっぺんにデプロイをリクエストすることができる
    - docker-composeのようなインターフェイスを採用
    - `variable` を使って変数解決を行う
    - 完全な情報を受け取る
        - 何かアタッチメントを追加したいときもすべてのアタッチメントを書く
          - ただしIDのみで良い
- 1つのリクエストで可能な限り多くの情報を受けとることで、高度なスケジューリングを行うことができる
- `/api/graph` がクエリ、 `/api/spec` がコマンド・ESである
    - つまり `/api/spec` を追っていれば操作がすべて追える
- annotationsはschedulerへのヒント

## 考察

### メリット

- `/api/spec` によってリソースの展開が非常に容易
- `/api/spec` を見ることですべての操作について確認を行うことができる
- `/api/graph` で直感的に現在のリソースを取得できる

### デメリット

- 値の優先度が分かりにくい
    - テンプレートを展開した時にどちらの値を優先するのか
- yamlの記法が複雑になりやすいので、学習コストが発生する
- relationsの記述が冗長
  - propertyはデフォルト値を柔軟に設定するしかないか？

## 懸念点

- `/api/spec` はスケールするのか
- 結局いつどうやってspecをデータストアに保存するのか
    - 削除はどうするのか -> ここでする議論ではない
- 値を変更したい時にその値だけを記述すると事を許すか
  - 要素を追加するのか、変更するのかがわかりにくい
  - べき等なことを売りにするならちゃんと入力してもらったほうが良いかもしれない

## 動作例

### `/api/spec`

- docker_composeのインターフェイスを参考にした
- 新規作成の際には `web` や `db` みたいな変数名を入れられる

1. request

- in this case, created volume has `a8d1d875-240a-445f-a569-10e00122e65b`.

```yaml
version: 0
annotations:
  n0stack.io/scheduling/same_host: true
  n0stack.io/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  web:
    type: resource/vm/kvm
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_network
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
      - object:
          variable: var_volume
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          id: a8d1d875-240a-445f-a569-10e00122e65b
          size: 100gb
        property:
          r: n0stack.jp/n0core/vm/attachments
  var_volume:
    type: resource/volume/nfs
    size: 10gb
  var_network:
    type: resource/network/vlan
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        enable_dhcp: true
        allocation_pool: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway_ip: 192.168.0.254
    parameters:
      id: 100
```

2. response

```yaml
id: 100
created_at: 1990/1/1 00:00
version: 0
annotations:
  n0stack.io/scheduling/same_host: true
  n0stack.io/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  web:
    type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: started
    name: web
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
          status: attached
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_network
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: d99163ed-0093-40a0-a61b-365a1aece509
          status: claimed
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
      - object:
          variable: var_volume
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: a8d1d875-240a-445f-a569-10e00122e65b
          status: claimed
          name: created_volume
          size: 100gb
        property:
          r: n0stack.jp/n0core/vm/attachments
  var_volume:
    type: resource/volume/nfs
    id: 2282dcee-d49f-4a6a-8a41-70e3e59a80cd
    status: claimed
    name: var_volume
    size: 10gb
  var_network:
    type: resource/network/vlan
    id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
    state: applied
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        enable_dhcp: true
        allocation_pool: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway_ip: 192.168.0.254
    parameters:
      id: 100
```

### `/api/spec` テンプレート機能

- テンプレートを変数として定義できる
    - 引数を与えることができる
    - WIP: DBにも保存できるようにするのか?
    - WIP: template_variableを変数として渡すか
- テンプレートを適用した上から要素を上書きすることができる

1. request

```yaml
version: 0
annotations:
  n0stack.io/scheduling/same_host: true
  n0stack.io/scheduling/host_id: hogehoge
spec:
  web-template:
    args:
      ip_octet: int
    type: template/resource/vm
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          ip_addresses:
            - 192.168.0.__ip_octet__
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
  web-0:
    type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: started
    template_variable:
      - web-template:
          ip_octet: 1
  web-1:
    type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c262
    status: started
    memory: 8gb
    template_variable:
      - web-template:
          ip_octet: 1
```

2. response

```yaml
version: 0
annotations:
  n0stack.io/scheduling/same_host: true
  n0stack.io/scheduling/host_id: hogehoge
spec:
  web-template:
    args:
      ip_octet: int
    type: template/resource/vm
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          ip_addresses:
            - 192.168.0.__ip_octet__
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: d99163ed-0093-40a0-a61b-365a1aece509
          status: claimed
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
  web-0:
    type: resource/vm/kvm
    template_variable:
      - web-template:
          ip_octet: 1
    relations:
      - object:
          type: resource/port
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
          status: attached
          ip_addresses:
            - 192.168.0.__ip_octet__
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: a8d1d875-240a-445f-a569-10e00122e65b
          status: claimed
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
  web-1:
    type: resource/vm
    memory: 8gb
    template_variable:
      - web-template:
          ip_octet: 1
    relations:
      - object:
          type: resource/port
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f82
          status: attached
          ip_addresses:
            - 192.168.0.__ip_octet__
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: a8d1d875-240a-445f-a569-10e00122e65d
          status: claimed
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
```

---

以下具体的な話

## エンドポイント

### `/api/spec`

#### `POST /api/spec`

- `/api/spec?version_from=100` のようなシーク機能をつける
- とあるバージョンのリソースを再現、ロールバックできるようなインターフェイスがほしい
  - `/api/spec?version=100&resource_id=8e3928cc-2afd-4bb4-beb4-6faeddb99a83` みたいな
- defaultの値とは
  - `name`: 変数名
  - `state`: 各リソースによる
- 変数代入、ID代入


```yaml
version: 0
spec:
  web:
    type: resource/vm/kvm
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
      - object:
          variable: var_volume
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          id: a8d1d875-240a-445f-a569-10e00122e65b
          size: 100gb
        property:
          r: n0stack.jp/n0core/vm/attachments
  var_volume:
    type: resource/volume/nfs
    size: 10gb
  var_network:
    type: resource/network/vlan
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        enable_dhcp: true
        allocation_pool: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway_ip: 192.168.0.254
    parameters:
      id: 100
```

#### response

- idを付与
- デフォルトの値を挿入
- 変数名は維持する

```yaml
id: 100
created_at: 1990/1/1 00:00
version: 0
spec:
  web:
    type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: started
    name: web
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    relations:
      - object:
          type: resource/port
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
          status: attached
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_natwork
              property:
                r: n0stack.jp/n0core/port/network
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: d99163ed-0093-40a0-a61b-365a1aece509
          status: claimed
          name: new_volume
          size: 10gb
        property:
          r: n0stack.jp/n0core/vm/attachments
          n0stack.jp/n0core/vm/boot_prority: 1
      - object:
          variable: var_volume
        property:
          r: n0stack.jp/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: a8d1d875-240a-445f-a569-10e00122e65b
          status: claimed
          name: created_volume
          size: 100gb
        property:
          r: n0stack.jp/n0core/vm/attachments
  var_volume:
    type: resource/volume/nfs
    id: 2282dcee-d49f-4a6a-8a41-70e3e59a80cd
    status: claimed
    name: var_network
    size: 10gb
  var_network:
    type: resource/network/vlan
    id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
    state: applied
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        enable_dhcp: true
        allocation_pool: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway_ip: 192.168.0.254
    parameters:
      id: 100
```

### `/api/graph`

- graphqlの予定

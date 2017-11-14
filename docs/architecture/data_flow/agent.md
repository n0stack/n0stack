# Scheduler => agent => Conductor (Messages)

## 要件

- リソースを間違いなく展開できるようにする必要がある
- 依存関係を `scheduler => agent => agent => conductor` のように効率的に解決できるようにする
    - WIP: 失敗した時のために次の候補も必要なのでは?
- メッセージを簡単にしたい
    - メンテナンス性向上

## 原則・モデル

- agentはステートレスであり、べき等なリソースを管理・抽象化するものである
- agentにはリソースのべき等な状態をリクエストする
- メッセージの中ではresource_idだけでリンクすることはない
    - データ構造をrelationsにそのまま入れる
- agentはリソースの依存関係を自分で解決できる
  - リソースの情報が足りない場合はリソースの作成プロセスをブロックし、ある一定期間自分に対するキューとして積んでおく

## 考察

### メリット

- メッセージからべき等であることを明示することでインフラ基盤としてわかりやすい
- リソースの種類を増やす際もデータ構造を1種類増やすだけで対応できる
- パラメータを増やす際も1つのprotobufファイルの操作で済む
- notifyでべき等な状態を広報することで依存しているものが依存されているものの情報を取得しやすい
    - VMはネットワークのブリッジを知る必要がある
    - agentがどんな情報が必要なのかという情報を考える必要がなくなる

### デメリット

- べき等なものとして定義できないものができたときにどうするのか
    - POSTのような機能が必要になった時にどうするのか
    - 例外を増やしてしまうのか

## 懸念点

- agentにメッセージを送っているときに新しいspecが来たらどうやって整合性を取るか
    - APIの時点でべき等なspecをclientに送ってもらう
      - 現在のステータスは関係ない
- VolumeのようにPOSTすると保存先ができるようなものはべき等にしにくい
    - schedulerで保存先まで指定する
    - `type: nfs` とかを送るとurlを発行する
    - urlを不完全にする `file:///` って送るとディレクトリを保管する
    - volumer一つのエージェントにつき一つのタイプしか対応しないのであれば、リクエストを送るだけでURLを作成できる
        - つまり、volumer-nfsに送れば `nfs://$ip/hoge/hoge` っていうのは簡単かつ正しく生成できる
        - とりあえずこれを採用
- すべてのフィールドに常に値を入れるか
    - 入れる
- いい加減agentのjoin, leaveを考えたい
    - 検討中
- 依存関係の解決をするときに次のリソースの候補がダメなときにどうやってメッセージを送るのか
    - volumer[0] => compute[0] の時に compute[0] が死んでいた場合 compute[1]に送る必要があるが、現状難しい
    - エラーハンドリングの項で提案

## 動作例

1. Request (volume, vm, network)

- vmはurlがないのでvolume、bridgeがないのでnetworkのアタッチメントに失敗する
- `topic: `
    - 本当だったら `vm/anycast` とかに投げたいところ
        - そうすることで任意のagentが作成を受信することができるので拡張性を持たせやすい

```yaml
spec_id: 100
object:
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
        hw_addr: ffffffffffff
        ip_addresses:
          - 192.168.0.1
        relatioins:
          - object:
              type: resource/network/vlan
              id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
              state: created
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
              relations:
                - property:
                    r: n0stack.jp/n0core/resource/scheduled
                  object:
                    type: agent/porter/vlan
                    id: a0c819fa-9dc2-4666-b7fd-d235a2551119
                    state: alived
                    host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
            property:
              r: n0stack.jp/n0core/port/network
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/vm/attachments
    - object:
        type: resource/volume/file
        id: d99163ed-0093-40a0-a61b-365a1aece509
        status: claimed
        name: new_volume
        size: 10gb
        relations:
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/volumer/file
              id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/vm/attachments
        n0stack.jp/n0core/vm/boot_prority: 1
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/compute/kvm
        id: 2463f81d-20d8-4395-a3c4-84a271a5b3a7
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

- `topic: `

```yaml
spec_id: 100
object:
  type: resource/volume/file
  id: d99163ed-0093-40a0-a61b-365a1aece509
  status: claimed
  name: new_volume
  size: 10gb
  relations:
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/volumer/file
        id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

- `topic: `

```yaml
spec_id: 100
object:
  type: resource/port
  id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
  status: attached
  hw_addr: ffffffffffff
  ip_addresses:
    - 192.168.0.1
  relatioins:
    - object:
        type: resource/network/vlan
        id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        state: created
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
        relations:
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/port/network
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

- `topic: `

```yaml
spec_id: 100
object:
  type: resource/network/vlan
  id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
  state: created
  name: vlan_network
  bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
  subnets:
    - cidr: 192.168.0.0/24
      enable_dhcp: true
      allocation_pool: 192.168.0.1-192.168.0.127
      nameservers:
        - 192.168.0.254
      gateway_ip: 192.168.0.254
  parameters:
    id: 100
  relations:
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

2. Notify (volume, network)

- networkが `parents` にも作成を通知することでportのブロック状態を解除する
- volumeがvmにできたことを通知するが、portができていないのでブロック状態を継続
- `topic: conductor/anycast`
- `topic: `

```yaml
spec_id: 100
msg: Succeeded to create resource.
level: SUCCESS
object:
  type: resource/volume/file
  id: d99163ed-0093-40a0-a61b-365a1aece509
  status: claimed
  name: new_volume
  url: file:///data/d99163ed-0093-40a0-a61b-365a1aece509
  size: 10gb
  relations:
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/volumer/file
        id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

- `topic: conductor/anycast`
- `topic: `

```yaml
spec_id: 100
msg: Succeeded to create resource.
level: SUCCESS
object:
  type: resource/network/vlan
  id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
  state: created
  name: vlan_network
  bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
  subnets:
    - cidr: 192.168.0.0/24
      enable_dhcp: true
      allocation_pool: 192.168.0.1-192.168.0.127
      nameservers:
        - 192.168.0.254
      gateway_ip: 192.168.0.254
  parameters:
    id: 100
  relations:
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

3. Notify (port)

- portが `parents` にも作成を通知することでvmのブロック状態を解除する
- `topic: conductor/anycast`
- `topic: `

```yaml
spec_id: 100
msg: Succeeded to create resource after waiting network resource(8451da31-5e3a-4c46-aa3a-2a557382a6cd).
level: SUCCESS
object:
  type: resource/port
  id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
  status: attached
  hw_addr: ffffffffffff
  ip_addresses:
    - 192.168.0.1
  relatioins:
    - object:
        type: resource/network/vlan
        id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        state: created
        name: vlan_network
        bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
        subnets:
          - cidr: 192.168.0.0/24
            enable_dhcp: true
            allocation_pool: 192.168.0.1-192.168.0.127
            nameservers:
              - 192.168.0.254
            gateway_ip: 192.168.0.254
        parameters:
          id: 100
        relations:
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/port/network
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

4. Notify (vm)

- `topic: conductor/anycast`

```yaml
spec_id: 100
msg: Succeeded to create resource after waiting network resource(8451da31-5e3a-4c46-aa3a-2a557382a6cd).
level: SUCCESS
object:
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
        hw_addr: ffffffffffff
        ip_addresses:
          - 192.168.0.1
        relatioins:
          - object:
              type: resource/network/vlan
              id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
              state: created
              name: vlan_network
              bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
              subnets:
                - cidr: 192.168.0.0/24
                  enable_dhcp: true
                  allocation_pool: 192.168.0.1-192.168.0.127
                  nameservers:
                    - 192.168.0.254
                  gateway_ip: 192.168.0.254
              parameters:
                id: 100
              relations:
                - property:
                    r: n0stack.jp/n0core/resource/scheduled
                  object:
                    type: agent/porter/vlan
                    id: a0c819fa-9dc2-4666-b7fd-d235a2551119
                    state: alived
                    host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
            property:
              r: n0stack.jp/n0core/port/network
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/vm/attachments
    - object:
        type: resource/volume/file
        id: d99163ed-0093-40a0-a61b-365a1aece509
        status: claimed
        name: new_volume
        url: file:///data/d99163ed-0093-40a0-a61b-365a1aece509
        size: 10gb
        relations:
          - property:
              r: n0stack.jp/n0core/resource/scheduled
            object:
              type: agent/volumer/file
              id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack.jp/n0core/vm/attachments
        n0stack.jp/n0core/vm/boot_prority: 1
    - property:
        r: n0stack.jp/n0core/resource/scheduled
      object:
        type: agent/compute/kvm
        id: 2463f81d-20d8-4395-a3c4-84a271a5b3a7
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

### エラーハンドリング

- エラーが発生した場合には Conductor がそれを回収し、Schedulerに以下のようなメッセージを送ることでスケジューリングと展開を再度行う
  - agentが返答を全くしないようなエラーをどうするか
    - タイムアウトを設定？
  - specを記述するか
    - 冗長だししなくても良さそう

```yaml
id: hosdfi  # エラーを受け取ったspec_id
annotations:
  n0stack.jp/n0core/scheduler/rescheduling: with_error
```

---

以下、具体的な設計

## Resources

- データ構造は [Resources](../data_structure/resources.md) と [Messsages](messages.md) を参照のこと

### VM

- TODO: スナップショットの管理の方法を考える
- TODO: CloudInitについて考える

### Snapshot

`type: snapshot`

### Volume

- requestのときはurlが空

### Network

- requestのときはbridgeが空

### Port

- WIP: `resource/port/floating_ip` の扱いについて

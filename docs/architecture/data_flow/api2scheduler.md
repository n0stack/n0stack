# API => Scheduler (Messages)

## 要件

- 複数のリソースをリクエストことができる
- scheduling hint(annotations)を受け取れる必要がある

## 原則・モデル

- あくまで複数のリソース・オブジェクトをリストとしてリクエストする機構
  - リストなのでdependenciesにはidのみを入力する
- 1つのリクエストで1つのスケジューリング単位
- specはagentへの命令、annotationはschedulerへの命令

## 考察

### メリット

- リソースのリクエストだけではなく、agentのJoin処理なども受け取れるように抽象化されている

### デメリット

## 懸念点

- apiで変数の解決からデータのパース、穴埋めまで行う必要があるので負荷が集中してしまうのではないか
- 何がvalidなリクエストなのか判断しにくい
    - schedulerが何も考えずに送ってしまうか
    - protobufに制約を加えることはできない

## 動作例

- [この](client2api.md)リクエストを受けた時の想定

```yaml
id: 100
annotations:
  n0stack/scheduling/same_host: true
  n0stack/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  - type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: running
    name: web
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    dependencies:
      - object:
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
        property:
          r: n0stack/n0core/vm/attachments
      - object:
          id: d99163ed-0093-40a0-a61b-365a1aece509
        property:
          r: n0stack/n0core/vm/attachments
          n0stack/n0core/vm/boot_prority: 1
      - object:
          id: 2282dcee-d49f-4a6a-8a41-70e3e59a80cd
        property:
          r: n0stack/n0core/vm/attachments
      - object:
          id: a8d1d875-240a-445f-a569-10e00122e65b
        property:
          r: n0stack/n0core/vm/attachments
  - type: resource/volume/nfs
    id: 2282dcee-d49f-4a6a-8a41-70e3e59a80cd
    status: allocated
    name: var_volume
    size: 10gb
  - type: resource/volume/file
    id: a8d1d875-240a-445f-a569-10e00122e65b
    status: allocated
    name: created_volume
    size: 100gb
  - type: resource/volume/file
    id: d99163ed-0093-40a0-a61b-365a1aece509
    status: allocated
    name: new_volume
    size: 10gb
  - type: resource/nic
    id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
    status: attached
    ip_addresses:
      - 192.168.0.1
    relatioins:
      - object:
          id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        property:
          r: n0stack/n0core/port/network
  - type: resource/network/vlan
    id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
    state: applied
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        dhcp:
          range: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway: 192.168.0.254
    parameters:
      vlan_id: 100
```

---

以下具体的な話

## message

## spec

- メッセージの構造については [ここ](messages.md) を参照のこと

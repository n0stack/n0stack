# Boot VirtualMachine with Image

In case of booting VirtualMachine `test` with Image `cloudimage-ubuntu-1804` on Network `test-network`.

## Example

```yaml
GenerateBlockStorage:
  type: Image
  action: GenerateBlockStorage
  args:
    image_name: cloudimage-ubuntu-1804
    tag: "latest"
    block_storage_name: test-blockstorage
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824
    limit_bytes: 10737418240

ApplyNetwork:
  type: Network
  action: ApplyNetwork
  args:
    name: test-network
    ipv4_cidr: 192.168.0.0/24
    annotations:
      n0core/provisioning/virtual_machine/vlan_id: "100"

CreateVirtualMachine:
  type: VirtualMachine
  action: CreateVirtualMachine
  args:
    name: test
    annotations:
      n0core/provisioning/virtual_machine/request_node_name: vm-host1
    request_cpu_milli_core: 10
    limit_cpu_milli_core: 1000
    request_memory_bytes: 536870912
    limit_memory_bytes: 536870912
    block_storage_names:
      - test-blockstorage
    nics:
      - network_name: debug-network
        ipv4_address: 192.168.0.1
    uuid: 056d2ccd-0c4c-44dc-a2c8-39a9d394b51f
    login_username: test
    ssh_authorized_keys:
      - ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey
  depends_on:
    - GenerateBlockStorage
    - ApplyNetwork
```

```sh
n0cli --api-endpoint=$api_ip:20180 do $path_of_previous_yaml
```


## Overview

- Image から BlockStorage を生成する
    - Image は Docker Image のように名前とタグによって管理されているため、タグを指定する必要がある
        - タグは `n0cli get image cloudimage-ubuntu-1804` の `tags` で確認することができる
    - `block_storage_name` で生成する BlockStorage の名前を指定する
        - VirtualMachine 生成時にVMとブロックストレージを接続するために用いる
    - まだスケジューリングに対応していないため、`annotations` の `n0core/provisioning/block_storage/request_node_name` で BlockStorage をどこのノードに配置するかを決める
        - ノードの名前は `n0cli get node` で確認できる
    - 生成する BlockStorage の容量は `10 GB (10737418240 Bytes)`
        - ゲストOSからはブロックストレージがこのサイズに見える
    - 生成する BlockStorage の実際に使う可能性のある容量は `1 GB (1073741824 Bytes)`
        - この値はスケジューリングなどに用いられる
- Network を作成 / 更新する
- VirtualMachineを作成する
    - `request_cpu_milli_core` で実際に使うであろうCPUコアを選択し、`limit_cpu_milli_core`で上限を指定する
        - `limit_cpu_milli_core` はCPUコア数を指定するため、 `limit_cpu_milli_core % 1000 == 0` である必要がある
        - この場合1コアのVMがたつ
    - `request_memory_bytes == limit_memory_bytes` である必要がある
        - この場合メモリ `512 MB (536870912 Bytes)`のVMがたつ
        - KVMのmemory ballooningは性能劣化が激しかったので、無効化しているため
    - まだスケジューリングに対応していないため、`annotations` の `n0core/provisioning/virtual_machine/request_node_name` で BlockStorage をどこのノードに配置するかを決める
    - `block_storage_names` で接続する BlockStorageを指定する
        - この場合、Image から作成した BlockStorage を接続している
    - `nics` でどの Network に接続するか指定する
        - この場合、作成した Network に `192.168.0.1` で接続することを宣言している
    - `uuid` は `uuidgen` などで適宜生成すること
    - 使っているゲストOSイメージが cloud-init に対応していた場合、`nics`で指定したIP、`login_username`で指定したユーザ、`ssh_authorized_keys`で指定したSSH公開鍵が設定される

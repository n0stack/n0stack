# Compute

## Principle

- Agentにはシンプルな命令方式のgRPCインターフェイス( kvm, iproute2 )を実装
- `ApplyCompute` で適用
- ステートについてはgoroutineでqmpのイベントを監視し、Componentで更新
  - `annotations.n0core/compute/state` にステートをいれる

## Example

```sh
grpc_cli call localhost:20180 n0stack.provisioning.ComputeService/ApplyCompute \
'metadata {
  name: "test-compute"
  version: 0
}
spec {
  vcpus: 1
  memory_bytes: 1073741824
  volumes {
    key: "test-volume"
    value: {
      volume_name: "test-volume"
    }
  }
  nics {
    key: "test-nic"
    value: {
      network_name: "test-network"
      hardware_address: "52:54:00:00:00:01"
      ip_addresses: "192.168.0.1"
    }
  }
}'
```

```sh
grpc_cli call localhost:20180 n0stack.provisioning.ComputeService/DeleteCompute \
'name: "test-compute"'
```

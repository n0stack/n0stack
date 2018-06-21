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
  annotations {
    key: "n0core/node_name"
    value: "test-node"
  }
  annotations {
    key: "n0core/qmp_path"
    value: "/var/lib/n0core/kvm/test-compute/monitor.sock"
  }
  annotations {
    key: "n0core/vnc_websocket_port"
    value: "59802"
  }
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

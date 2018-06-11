# Volume

## Principle

- Agentにはシンプルな命令方式のgRPCインターフェイス( github.com/n0stack/n0core/provisioning/node/qcow2 )を実装
- APIを叩かれた時に必要に応じてAgentのgRPCをたたく

## Example

```sh
% grpc_cli call localhost:20180 n0stack.provisioning.VolumeService/ApplyVolume 'metadata {
  name: "test-node"
  annotations {
    key: "n0core/url"
    value: "file:///tmp/test.qcow2"
  }
  annotations {
    key: "n0core/node_name"
    value: "test-node"
  }
  version: 0
}
spec {
  bytes: 1073741824
}'
connecting to localhost:20180
metadata {
  name: "test-node"
  annotations {
    key: "n0core/node_name"
    value: "test-node"
  }
  annotations {
    key: "n0core/url"
    value: "file:///tmp/test.qcow2"
  }
  version: 1
}
spec {
  bytes: 1073741824
}
status {
  state: AVAILABLE
}

Rpc succeeded with OK status
```

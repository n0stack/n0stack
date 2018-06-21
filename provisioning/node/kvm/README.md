# KVM

- runlevelについてはgoroutineからAPIに通知を行う
- 変数にデフォルト値はセットしない

## Example

```sh
grpc_cli call localhost:20181 n0stack.n0core.kvm.KVMService/ApplyKVM \
'kvm {
  uuid: "a425dede-b6b3-5572-a3e2-7de689a6d8a5"
  name: "test-compute"
  vnc_websocket_port: 5000
  qmp_path: "/tmp/monitor.sock"

  cpu_cores: 1
  memory_bytes: 1073741824

  storages: {
    key: "test-volume"
    value: {
      url: "file:///tmp/ubuntu.qcow2"
      boot_index: 1
    }
  }

  nics: {
    key: "test-nic"
    value: {
      tap_name: "test-tap"
      hw_addr: "52:54:00:00:00:01"
    }
  }
}'
```

```sh
grpc_cli call localhost:20181 n0stack.n0core.kvm.KVMService/DeleteKVM \
'name: "test-vm"'
```

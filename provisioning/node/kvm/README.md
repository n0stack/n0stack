# KVM

- runlevelについてはgoroutineからAPIに通知を行う

## Example

```sh
grpc_cli call localhost:20181 n0stack.n0core.kvm.KVMService/ApplyKVM \
'kvm {
  uuid: "9f8f7a4e-d314-4135-bebc-e0a44e7bcbe9"
  name: "test-vm"
  vnc_websocket_port: 5000
  qmp_path: "/tmp/monitor.sock"

  cpu_cores: 1
  memory_bytes: 1073741824

  volumes: {
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

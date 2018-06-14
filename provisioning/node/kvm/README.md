# KVM

- runlevelについてはgoroutineからAPIに通知を行う

## Example

```sh
grpc_cli call localhost:20181 n0stack.n0core.iproute2.Iproute2Service/ApplyTap \
'kvm {
  uuid: "9f8f7a4e-d314-4135-bebc-e0a44e7bcbe9"
  name: "test-tap"
  cpu_cores: 1
  memory_bytes: 1073741824
  qmp_path: "/tmp/monitor.sock"
}'
```

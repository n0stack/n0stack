# iproute2

## Example

```sh
grpc_cli call localhost:20181 n0stack.n0core.iproute2.Iproute2Service/ApplyTap \
'tap {
  name: "test-tap"
  bridge_name: "test-bridge"
  type: FLAT
}'
```

```sh
grpc_cli call localhost:20181 n0stack.n0core.iproute2.Iproute2Service/DeleteTap \
'name: "test-tap"'
```

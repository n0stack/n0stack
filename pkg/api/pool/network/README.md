# Network

## Example

```
grpc_cli call localhost:20182 n0stack.pool.NetworkService/ListNetworks ''
```

```
grpc_cli call localhost:20182 n0stack.pool.NetworkService/GetNetwork '
name: "test-network"
'
```

```
grpc_cli call localhost:20182 n0stack.pool.NetworkService/ApplyNetwork '\
name: "test-network"
ipv4_cidr: "10.100.100.0/24"
domain: "test.local"
'
```

```
grpc_cli call localhost:20182 n0stack.pool.NetworkService/ReserveNetworkInterface '
name: "test-network"
network_interface_name: "test-reserve"
'
```

```
grpc_cli call localhost:20182 n0stack.pool.NetworkService/ReleaseNetworkInterface '
name: "test-network"
network_interface_name: "test-reserve"
'
```

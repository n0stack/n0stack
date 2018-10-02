# Provisioning

## Example

### BlockStorage

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/CreateEmptyBlockStorage '\
name: "test-empty-volume"
annotations {
  key: "n0core/provisioning/request_node_name"
  value: "test"
}
request_bytes: 1024
limit_bytes: 1073741824
'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/CreateBlockStorageWithDownloading '\
name: "test-ubuntu-volume"
annotations {
  key: "n0core/provisioning/request_node_name"
  value: "test"
}
request_bytes: 1073741824
limit_bytes: 10737418240
source_url: "http://cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img"
'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/ListBlockStorages ''
```

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/GetBlockStorage 'name: "test-ubuntu-volume"'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/SetAvailableBlockStorage 'name: test-ubuntu-volume"'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.BlockStorageService/SetInuseBlockStorage 'name: "test-ubuntu-volume"'
```


### Virtual machine

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/CreateVirtualMachine '\
name: "test-vm"
annotations {
  key: "n0core/provisioning/request_node_name"
  value: "test"
}
request_cpu_milli_core: 10
limit_cpu_milli_core: 1000

request_memory_bytes: 1073741824
limit_memory_bytes: 1073741824

block_storage_names: "test-ubuntu-volume"

nics {
  network_name: "test-network"
}
'
```

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/ListVirtualMachines ''
```

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/GetVirtualMachine 'name: "test-vm"'
```

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/BootVirtualMachine 'name: "test-vm"'
```

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/RebootVirtualMachine '
name: "test-vm"
hard: true
'
```

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/ShutdownVirtualMachine '
name: "test-vm"
hard: true
'
```

# Provisioning

## Example

### Volume

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/CreateEmptyVolume \
'metadata {
  name: "test-empty-volume"
  annotations {
    key: "n0core/provisioning/request_node_name"
    value: "test"
  }
}
spec {
  request_bytes: 1024
  limit_bytes: 1073741824
}'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/CreateVolumeWithDownloading \
'metadata {
  name: "test-ubuntu-volume"
  annotations {
    key: "n0core/provisioning/request_node_name"
    value: "test"
  }
}
spec {
  request_bytes: 1073741824
  limit_bytes: 10737418240
}
source_url: "http://cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img"
'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/ListVolumes ''
```

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/GetVolume 'name: "test-ubuntu-volume"'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/SetAvailableVolume 'name: test-ubuntu-volume"'
```

```
grpc_cli call localhost:20183 n0stack.provisioning.VolumeService/SetInuseVolume 'name: "test-ubuntu-volume"'
```


### Virtual machine

```
grpc_cli call localhost:20184 n0stack.provisioning.VirtualMachineService/CreateVirtualMachine \
'metadata {
  name: "test-vm"
  annotations {
    key: "n0core/provisioning/request_node_name"
    value: "test"
  }
}
spec {
  request_cpu_milli_core: 10
  limit_cpu_milli_core: 1000

  request_memory_bytes: 1073741824
  limit_memory_bytes: 1073741824

  volume_names: "test-ubuntu-volume"

  nics {
    network_name: "test-network"
  }
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

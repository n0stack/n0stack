# Boot VirtualMachine with ISO

## Fetch and register Ubuntu 18.04 Cloud Images
```yaml
FetchBlockStorage:
  type: BlockStorage
  action: FetchBlockStorage
  args:
    name: cloudimage-ubuntu-1804
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824 # 1GiB
    limit_bytes: 10737418240 # 10GiB
    source_url: https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img
  ignore_error: true

ApplyImage:
  type: Image
  action: ApplyImage
  args:
    name: cloudimage-ubuntu

RegisterBlockStorage:
  type: Image
  action: RegisterBlockStorage
  args:
    image_name: cloudimage-ubuntu
    block_storage_name: cloudimage-ubuntu-1804
    tags:
      - "latest"
      - "1804"
  depends_on:
    - FetchBlockStorage
    - ApplyImage
```

## Create VM using registered image

```yaml
GenerateBlockStorage:
  type: Image
  action: GenerateBlockStorage
  args:
    image_name: cloudimage-ubuntu
    block_storage_name: test-with-iso
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824 # 1GiB
    limit_bytes: 10737418240 # 10GiB
    tag: "1804"
  depends_on:
    - RegisterBlockStorage

ApplyNetwork:
  type: Network
  action: ApplyNetwork
  args:
    name: test-network
    ipv4_cidr: 192.168.0.0/24
    annotations:
      n0core/provisioning/virtual_machine/vlan_id: "100"

CreateVirtualMachine:
  type: VirtualMachine
  action: CreateVirtualMachine
  args:
    name: test-with-iso
    annotations:
      n0core/provisioning/virtual_machine/request_node_name: vm-host1
    request_cpu_milli_core: 10
    limit_cpu_milli_core: 1000
    request_memory_bytes: 1073741824 # 1GiB
    limit_memory_bytes: 1073741824 # 1GiB
    block_storage_names:
      - test-with-iso
    nics:
      - network_name: test-network
    login_username: ubuntu
    ssh_authorized_keys:
      # - ecdsa-sha2-nistp256 ...
  depends_on:
    - GenerateBlockStorage

OpenConsole:
  type: VirtualMachine
  action: OpenConsole
  args:
    name: test-with-iso
  depends_on:
    - CreateVirtualMachine
```
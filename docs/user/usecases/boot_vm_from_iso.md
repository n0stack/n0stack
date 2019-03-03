# Boot VirtualMachine from ISO

## Fetch and register Ubuntu 18.04 Cloud Images

```yaml
FetchISO:
  type: BlockStorage
  action: FetchBlockStorage
  args:
    name: cloudimage-ubuntu-1804
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824 # 1GiB
    limit_bytes: 10737418240 # 10GiB
    source_url: https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img

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
    name: test-vm
    annotations:
      n0core/provisioning/virtual_machine/request_node_name: vm-host1
    request_cpu_milli_core: 10
    limit_cpu_milli_core: 1000
    request_memory_bytes: 1073741824 # 1GiB
    limit_memory_bytes: 1073741824 # 1GiB
    block_storage_names:
      - cloudimage-ubuntu-1804
    nics:
      - network_name: test-network
        ipv4_address: 192.168.0.1
    # cloud-config related options:
    login_username: n0user
    ssh_authorized_keys:
      - ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey
  depends_on:
    - CreateBlockStorage

# You need to set password for user to login via console (not set if default)
OpenConsole:
  type: VirtualMachine
  action: OpenConsole
  args:
    name: test-vm
  depends_on:
    - CreateVirtualMachine
```

Then, you can login virtual machine via ssh by `n0user` user using key below:

```
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIBAQh+adEg/rjqj9qLE0jI4EqV8kZFDzWTASAwvx6HWdoAoGCCqGSM49
AwEHoUQDQgAEhOjA+fY6XV4K9c3ldX4tvqN9fOANtfIR21rJp0NQm0Wtw3abaaML
UHbRUECglxm1JiSaOuWVLTpDbpN7mxNi8Q==
-----END EC PRIVATE KEY-----
```

(Ubuntu 18.04 Cloud Image doesn't allow password login to ssh configured above, so you need set password if need to access via VNC console)

## Inverse action

```yaml
Delete_test-vm:
  type: VirtualMachine
  action: DeleteVirtualMachine
  args:
    name: test-vm

Delete_blockstorage:
  type: BlockStorage
  action: DeleteBlockStorage
  args:
    name: cloudimage-ubuntu-1804
  depends_on:
    - Delete_test-vm

Delete_test-network:
  type: Network
  action: DeleteNetwork
  args:
    name: test-network
  depends_on:
    - Delete_test-vm
```

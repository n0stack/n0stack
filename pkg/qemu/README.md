# QEMU

## Features

- CPU は `host` で渡している
- Ballooning はしない
- 基本的にはvirtioで接続
- SCSIコントローラを作成している

## Dependency packages

- qemu-kvm
<!-- - ovmf -->

```sh
apt install -y \
    qemu-kvm
```

## Test parameters

### DISABLE_KVM

When setting environment variable, enable KVM.

```sh
DISABLE_KVM=1 make test-small
```

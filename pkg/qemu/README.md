# QEMU

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

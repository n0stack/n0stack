# Resources

One of [Objects](objects.md).

**If you want to get latest informations, please see `/proto`.**

## VM

```yaml
id: 13bae4ae-67f3-456a-ab05-a217d7cf0861
type: resource/vm/kvm
name: hogehoge
state: started
arch: amd64
vcpus: 2
memory: 4gb
vnc_password: hogehoge
relations:
  - object:
      id: 0a0615bf-8d26-4e9f-bfbc-bbd0890fcd4f
      type: resource/port
      name: port
      state: attached
      hw_addr: ffffffffffff
      ip_addrs:
        - 192.168.0.1
        - fe08::1
      relations:
        - object:
            id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
            type: resource/network/vlan
            name: hogehoge
            state: applied
            bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
            subnets:
              - cidr: 192.168.0.0/24
                enable_dhcp: true
                allocation_pool: 192.168.0.1-192.168.0.127
                nameservers:
                  - 192.168.0.254
                gateway_ip: 192.168.0.254
            parameters:
              id: 100
          property:
            r: n0stack.jp/n0core/port/network
          direction: ->
    property:
      r: n0stack.jp/n0core/vm/attachments
    direction: ->
  - object:
      type: resource/volume/file
      id: 486274b2-49e4-4bcd-a60d-4f627ce8c041
      state: claimed
      name: hogehoge
      size: 10gb
      url: file:///data/hoge
    property:
      r: n0stack.jp/n0core/vm/attachments
      n0stack.jp/n0core/vm/boot_priority: 1
    direction: ->
```

### state

- started
- stopped
- restarted
- halted
- destroyed

## Volume

```yaml
type: resource/volume/file
id: 486274b2-49e4-4bcd-a60d-4f627ce8c041
state: claimed
name: hogehoge
size: 10gb
url: file:///data/hoge
```

### state

- claimed
- released
- destroyed

## Network

```yaml
id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
type: resource/network/vlan
name: hogehoge
state: applied
bridge: br-flat
subnets:
  - cidr: 192.168.0.0/24
    enable_dhcp: true
    allocation_pool: 192.168.0.1-192.168.0.127
    nameservers:
      - 192.168.0.254
    gateway_ip: 192.168.0.254
parameters:
  id: 100
```

### state

- applied
- halted

## Port

```yaml
id: 0a0615bf-8d26-4e9f-bfbc-bbd0890fcd4f
type: resource/port
name: port
state: attached
hw_addr: ffffffffffff
ip_addrs:
  - 192.168.0.1
  - fe08::1
relations:
  - object:
      id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
      type: resource/network/vlan
      name: hogehoge
      state: applied
      bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
      subnets:
        - cidr: 192.168.0.0/24
          enable_dhcp: true
          allocation_pool: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway_ip: 192.168.0.254
      parameters:
        id: 100
    property:
      r: n0stack.jp/n0core/port/network
    direction: ->
```

### state

- attached
- detached

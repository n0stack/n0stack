# Example of launching VM

## 1. End users send to API (`/api/spec/`).

See details in [here](data_flow/client2api.md).

### 1.1 request

```yaml
version: 0
annotations:
  n0stack/scheduling/same_host: true
  n0stack/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  web:
    type: resource/vm/kvm
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    dependencies:
      - object:
          type: resource/nic
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_network
              property:
                r: n0stack/n0core/port/network
        property:
          r: n0stack/n0core/vm/attachments
      - object:
          type: resource/volume/file
          name: new_volume
          size: 10gb
        property:
          r: n0stack/n0core/vm/attachments
          n0stack/n0core/vm/boot_prority: 1
  var_network:
    type: resource/network/vlan
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        dhcp:
          range: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway: 192.168.0.254
    parameters:
      vlan_id: 100
```

### 1.2 response

```yaml
id: 100
created_at: 1990/1/1 00:00
version: 0
annotations:
  n0stack/scheduling/same_host: true
  n0stack/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  web:
    type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: running
    name: web
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    dependencies:
      - object:
          type: resource/nic
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
          status: attached
          ip_addresses:
            - 192.168.0.1
          relatioins:
            - object:
                variable: var_network
              property:
                r: n0stack/n0core/port/network
        property:
          r: n0stack/n0core/vm/attachments
      - object:
          type: resource/volume/file
          id: d99163ed-0093-40a0-a61b-365a1aece509
          status: allocated
          name: new_volume
          size: 10gb
        property:
          r: n0stack/n0core/vm/attachments
          n0stack/n0core/vm/boot_prority: 1
  var_network:
    type: resource/network/vlan
    id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
    state: applied
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        dhcp:
          range: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway: 192.168.0.254
    parameters:
      vlan_id: 100
```

## 2. APIs proccess.

- Store version record with spec.

## 3. APIs send message `Spec` to schedulers (`topic: scheduler`).

See details in [here](data_flow/api2scheduler.md).

```yaml
id: 100
annotations:
  n0stack/scheduling/same_host: true
  n0stack/scheduling/host_id: 85cf3a3a-18e4-4fe2-b406-9e79079cae07
spec:
  - type: resource/vm/kvm
    id: 56410722-d507-472a-a800-c89211b7c261
    status: running
    name: web
    arch: amd64
    vcpus: 2
    memory: 4gb
    vnc_password: hogehoge
    dependencies:
      - object:
          id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
        property:
          r: n0stack/n0core/vm/attachments
      - object:
          id: d99163ed-0093-40a0-a61b-365a1aece509
        property:
          r: n0stack/n0core/vm/attachments
          n0stack/n0core/vm/boot_prority: 1
  - type: resource/volume/file
    id: d99163ed-0093-40a0-a61b-365a1aece509
    status: allocated
    name: new_volume
    size: 10gb
  - type: resource/nic
    id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
    status: attached
    ip_addresses:
      - 192.168.0.1
    relatioins:
      - object:
          id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        property:
          r: n0stack/n0core/port/network
  - type: resource/network/vlan
    id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
    state: applied
    name: vlan_network
    subnets:
      - cidr: 192.168.0.0/24
        dhcp:
          range: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway: 192.168.0.254
    parameters:
      vlan_id: 100
```

## 4. Schedulers proccess.

- Issue `SCHEDULED` event.
- Add the relation for agent to store result of sheduled.

## 5. Agents create resources.

See details in [here](data_flow/agent.md).

### 5.1 Schedulers send message `Request` to each agents.

- `topic: $agent_id`

```yaml
spec_id: 100
object:
  type: resource/vm/kvm
  id: 56410722-d507-472a-a800-c89211b7c261
  status: running
  name: web
  arch: amd64
  vcpus: 2
  memory: 4gb
  vnc_password: hogehoge
  dependencies:
    - object:
        type: resource/nic
        id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
        status: attached
        hw_addr: ffffffffffff
        ip_addresses:
          - 192.168.0.1
        relatioins:
          - object:
              type: resource/network/vlan
              id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
              state: applied
              name: vlan_network
              subnets:
                - cidr: 192.168.0.0/24
                  dhcp:
                    range: 192.168.0.1-192.168.0.127
                    nameservers:
                      - 192.168.0.254
                    gateway: 192.168.0.254
              parameters:
                vlan_id: 100
              dependencies:
                - property:
                    r: n0stack/n0core/resource/scheduled
                  object:
                    type: agent/porter/vlan
                    id: a0c819fa-9dc2-4666-b7fd-d235a2551119
                    state: alived
                    host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
            property:
              r: n0stack/n0core/port/network
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/vm/attachments
    - object:
        type: resource/volume/file
        id: d99163ed-0093-40a0-a61b-365a1aece509
        status: allocated
        name: new_volume
        size: 10gb
        dependencies:
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/volumer/file
              id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/vm/attachments
        n0stack/n0core/vm/boot_prority: 1
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/compute/kvm
        id: 2463f81d-20d8-4395-a3c4-84a271a5b3a7
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

- `topic: $agent_id`

```yaml
spec_id: 100
object:
  type: resource/volume/file
  id: d99163ed-0093-40a0-a61b-365a1aece509
  status: allocated
  name: new_volume
  size: 10gb
  dependencies:
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/volumer/file
        id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

- `topic: $agent_id`

```yaml
spec_id: 100
object:
  type: resource/nic
  id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
  status: attached
  hw_addr: ffffffffffff
  ip_addresses:
    - 192.168.0.1
  relatioins:
    - object:
        type: resource/network/vlan
        id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        state: applied
        name: vlan_network
        subnets:
          - cidr: 192.168.0.0/24
            dhcp:
              range: 192.168.0.1-192.168.0.127
              nameservers:
                - 192.168.0.254
              gateway: 192.168.0.254
        parameters:
          vlan_id: 100
        dependencies:
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/port/network
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

- `topic: $agent_id`

```yaml
spec_id: 100
object:
  type: resource/network/vlan
  id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
  state: applied
  name: vlan_network
  bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
  subnets:
    - cidr: 192.168.0.0/24
      dhcp:
        range: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway: 192.168.0.254
  parameters:
    vlan_id: 100
  dependencies:
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
depended_by:
```

### 5.2 Agents send to conductor and depended agent.

- `topic: conductor`
- `topic: $agent_id`

```yaml
spec_id: 100
msg: Succeeded to create resource.
level: SUCCESS
object:
  type: resource/volume/file
  id: d99163ed-0093-40a0-a61b-365a1aece509
  status: allocated
  name: new_volume
  url: file:///data/d99163ed-0093-40a0-a61b-365a1aece509
  size: 10gb
  dependencies:
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/volumer/file
        id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

- `topic: conductor`
- `topic: $agent_id`

```yaml
spec_id: 100
msg: Succeeded to create resource.
level: SUCCESS
object:
  type: resource/network/vlan
  id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
  state: applied
  name: vlan_network
  bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
  subnets:
    - cidr: 192.168.0.0/24
      dhcp:
        range: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway: 192.168.0.254
  parameters:
    vlan_id: 100
  dependencies:
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

### 5.3 Agents send to conductor and depended agent 2nd.

- `topic: conductor`
- `topic: $agent_id`

```yaml
spec_id: 100
msg: Succeeded to create resource after waiting network resource(8451da31-5e3a-4c46-aa3a-2a557382a6cd).
level: SUCCESS
object:
  type: resource/nic
  id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
  status: attached
  hw_addr: ffffffffffff
  ip_addresses:
    - 192.168.0.1
  relatioins:
    - object:
        type: resource/network/vlan
        id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
        state: applied
        name: vlan_network
        bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
        subnets:
          - cidr: 192.168.0.0/24
            dhcp:
              range: 192.168.0.1-192.168.0.127
              nameservers:
                - 192.168.0.254
              gateway: 192.168.0.254
        parameters:
          vlan_id: 100
        dependencies:
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/port/network
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/porter/vlan
        id: a0c819fa-9dc2-4666-b7fd-d235a2551119
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

### 5.3 Agents send to conductor.

- `topic: conductor`

```yaml
spec_id: 100
msg: Succeeded to create resource after waiting network resource(8451da31-5e3a-4c46-aa3a-2a557382a6cd).
level: SUCCESS
object:
  type: resource/vm/kvm
  id: 56410722-d507-472a-a800-c89211b7c261
  status: running
  name: web
  arch: amd64
  vcpus: 2
  memory: 4gb
  vnc_password: hogehoge
  dependencies:
    - object:
        type: resource/nic
        id: a1f6a79b-7ad0-499e-b912-44eb73ea0f81
        status: attached
        hw_addr: ffffffffffff
        ip_addresses:
          - 192.168.0.1
        relatioins:
          - object:
              type: resource/network/vlan
              id: 8451da31-5e3a-4c46-aa3a-2a557382a6cd
              state: applied
              name: vlan_network
              bridge: br-vlan-0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
              subnets:
                - cidr: 192.168.0.0/24
                  dhcp:
                    range: 192.168.0.1-192.168.0.127
                    nameservers:
                      - 192.168.0.254
                    gateway: 192.168.0.254
              parameters:
                vlan_id: 100
              dependencies:
                - property:
                    r: n0stack/n0core/resource/scheduled
                  object:
                    type: agent/porter/vlan
                    id: a0c819fa-9dc2-4666-b7fd-d235a2551119
                    state: alived
                    host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
            property:
              r: n0stack/n0core/port/network
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/porter/vlan
              id: a0c819fa-9dc2-4666-b7fd-d235a2551119
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/vm/attachments
    - object:
        type: resource/volume/file
        id: d99163ed-0093-40a0-a61b-365a1aece509
        status: allocated
        name: new_volume
        url: file:///data/d99163ed-0093-40a0-a61b-365a1aece509
        size: 10gb
        dependencies:
          - property:
              r: n0stack/n0core/resource/scheduled
            object:
              type: agent/volumer/file
              id: 264b01b8-aeb5-478a-ad11-715fbc86d2f6
              state: alived
              host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
      property:
        r: n0stack/n0core/vm/attachments
        n0stack/n0core/vm/boot_prority: 1
    - property:
        r: n0stack/n0core/resource/scheduled
      object:
        type: agent/compute/kvm
        id: 2463f81d-20d8-4395-a3c4-84a271a5b3a7
        state: alived
        host_id: 8bce7696-f641-411c-a0ce-6ed066d912a3
```

## 6. Conductors proccess.

- `APPLIED` event issue.
- Apply result to GraphDB.

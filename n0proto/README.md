# n0proto

Protobuf definitions for all of n0stack services.

## Overview

see also [docs](https://docs.n0st.ac/en/latest/user/overview_n0proto.html).

## How to build

- Required Docker
- Generating Golang and Python files on n0proto.go and n0proto.py

```
cd ..
make build-n0proto-on-docker
```

## Principles

- Do not define variables that change with implementation, such values ​​should be placed in "annotations".
    - e.g. VLAN ID and VXLAN ID

### Standard fields

- Metadata (1 ~ 9)
- Spec (10 ~ 49)
- Status (50 ~)

```pb
  // Name is a unique field.
  string name = 1;
  // string namespace = 2;

  // Annotations can store metadata used by the system for control.
  // In particular, implementation-dependent fields that can not be set as protobuf fields are targeted.
  // The control specified by n0stack may delete metadata specified by the user.
  map<string, string> annotations = 3;

  // Labels stores user-defined metadata.
  // The n0stack system must not rewrite this value.
  map<string, string> labels = 4;
```

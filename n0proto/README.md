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
  // Make unique!!
  string name = 1;
  // string namespace = 2;

  map<string, string> annotations = 3;
  // map<string, string> labels = 4;
```

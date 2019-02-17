# Quick Start

## n0cli

The n0cli is a CLI tool to call n0stack gRPC APIs.

### Installation

with docker

```sh
docker pull n0stack/n0stack
docker run -it --rm -v /usr/local/bin:/dst n0stack/n0stack cp /usr/local/bin/n0cli /dst/
```

### Usage

- See also command help.

```sh
$ n0cli --api-endpoint=$api_ip:20180 get node
{
  "nodes": [
    {
      "name": "vm-host1",
      "annotations": {
        "github.com/n0stack/n0stack/n0core/agent_version": "52"
      },
      "address": "192.168.122.10",
      "serial": "Specified",
      "cpu_milli_cores": 1000,
      "memory_bytes": "1033236480",
      "storage_bytes": "107374182400",
      "unit": 1,
      "state": "Ready",
      "reserved_computes": {
        "debug_ipv6": {
          "annotations": {
            "n0core/provisioning/virtual_machine/virtual_machine/reserved_by": "debug_ipv6"
          },
          "request_cpu_milli_core": 10,
          "limit_cpu_milli_core": 1000,
          "request_memory_bytes": "536870912",
          "limit_memory_bytes": "536870912"
        }
      },
      "reserved_storages": {
        "debug-ipv6-network": {
          "annotations": {
            "n0core/provisioning/block_storage/reserved_by": "debug-ipv6-network"
          },
          "request_bytes": "1073741824",
          "limit_bytes": "10737418240"
        },
        "debug_ipv6_network": {
          "annotations": {
            "n0core/provisioning/block_storage/reserved_by": "debug_ipv6_network"
          },
          "request_bytes": "1073741824",
          "limit_bytes": "10737418240"
        },
        "ubuntu-1804": {
          "annotations": {
            "n0core/provisioning/block_storage/reserved_by": "ubuntu-1804"
          },
          "request_bytes": "1073741824",
          "limit_bytes": "10737418240"
        }
      }
    }
  ]
}
```

### Examples

See also [Usecases](usecases/README.rst).

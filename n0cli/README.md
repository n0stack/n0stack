# n0ctl

CLI for end-user

## Command

## Example

```sh
% n0core/bin/n0ctl get network
get.go:122: [DEBUG] Connected to 'localhost:20180'
got error response: rpc error: code = NotFound desc =

% n0core/bin/n0ctl do n0core/examples/n0ctl/get.yaml
do.go:44: [DEBUG] Connected to 'localhost:20180'
---> Task 'g3' is started
dag.go:169: [DEBUG] Task 'g3' is started: &{ResourceType:Node Action:ListNodes Args:map[] DependOn:[] child:[g1 g2] depends:0}
---> [ 1/3 ] Task 'g3' is finished
--- Response ---
{
  "nodes": [
    {
      "name": "mock-node",
      "version": 10,
      "address": "10.20.180.4",
      "cpu_milli_cores": 4000,
      "memory_bytes": 16692797440,
      "storage_bytes": 107374182400,
      "state": 1
    }
  ]
}
---> Task 'g1' is started
dag.go:205: [DEBUG] Task 'g1' is started: &{ResourceType:Network Action:ApplyNetwork Args:map[ipv4_cidr:10.100.100.0/24 domain:test.local name:test-network] DependOn:[g3] child:[] depends:0}
---> Task 'g2' is started
dag.go:205: [DEBUG] Task 'g2' is started: &{ResourceType:Node Action:GetNode Args:map[name:mock-node] DependOn:[g3] child:[] depends:0}
---> [ 2/3 ] Task 'g2' is finished
--- Response ---
{
  "name": "mock-node",
  "version": 10,
  "address": "10.20.180.4",
  "cpu_milli_cores": 4000,
  "memory_bytes": 16692797440,
  "storage_bytes": 107374182400,
  "state": 1
}
---> [ 3/3 ] Task 'g1' is finished
--- Response ---
{
  "name": "test-network",
  "version": 1,
  "ipv4_cidr": "10.100.100.0/24",
  "domain": "test.local",
  "state": 2
}
DAG tasks are completed

% bin/n0ctl get network
get.go:122: [DEBUG] Connected to 'localhost:20180'
{
  "networks": [
    {
      "name": "test-network",
      "ipv4_cidr": "10.100.100.0/24",
      "domain": "test.local",
      "state": 2
    }
  ]
}
```

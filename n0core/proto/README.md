n0stack Protobuf Schemas
===

# Policy

- Each `.proto` file must define one message / enum
- `all.proto` imports all `.proto` file without itself

# How to use

- Compiled file (for python) is located under `python` directory
- `import` each file and use:

# Defining new `.proto` file

1. Create `.proto` file
1. Add `import "...";` to `all.proto`
1. Compile `protoc -I. --python_out=python new_file.proto all.proto` (Considering `pwd` is in this directory)

# Example usage in Python

```
>>> import vm_power_state_request_pb2
>>> import vm_power_state_pb2
>>> a = vm_power_state_request_pb2.VMPowerStateRequest(host="host", name="name", status=vm_power_state_pb2.POWEROFF)
>>> a.SerializeToString()
b'\n\x04host\x12\x04name\x18\x02'
>>> serialized = a.SerializeToString()
>>> deserialized = vm_power_state_request_pb2.VMPowerStateRequest()
>>> deserialized.ParseFromString(serialized)
>>> deserialized
host: "host"
name: "name"
status: POWEROFF
```
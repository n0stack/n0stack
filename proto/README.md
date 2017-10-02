n0stack Protobuf Schemas
===

# Policy

- Each `.proto` file must define one message / enum
- `N0stackMessage` defined in `n0stack_message.proto` is a superclass of all message types.

# How to use

- Compile `*.proto` file to Python classes, use `protoc` compiler
  - Use local compiler: `protoc -I. --python_out=../n0core/proto new_file.proto all.proto` 
  - Or `protoc` in Docker: `docker run -it --rm -v $PWD:/src:rw -v $PWD/../n0core/proto:/dst nanoservice/protobuf --python_out=/dst **/*.proto`
  - (Considering working in this directory)
- `import` each file and use:

# Defining new `.proto` file

1. Create `.proto` file (for example: `new_file.proto`)
1. Make it message under sub-message of `N0stackMessage`

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
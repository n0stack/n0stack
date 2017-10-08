n0stack Protobuf Schemas
===

# Policy

- Each `.proto` file must define one message / enum
- `N0stackMessage` defined in `n0stack_message.proto` is a superclass of all message types.

# How to use

- Compile `*.proto` file to Python classes, use `protoc` compiler
  - Use local compiler: `protoc -I. --python_out=../n0core/proto new_file.proto all.proto` 
  - Or `protoc` in Docker: `docker run -it --rm -v $PWD:/src:rw -v $PWD/..:/dst nanoservice/protobuf --python_out=/dst **/*.proto`
  - (Considering working in this directory (n0core/proto) )
- `import` each file and use:

# Defining new `.proto` file

1. Create `.proto` file (for example: `new_file.proto`)
1. Make it message under sub-message of `N0stackMessage`

# Example usage in Python

```python
from n0core.lib.proto import UpdateVMPowerStateRequest, VMPowerState

req = UpdateVMPowerStateRequest(id="some_vm",
                                status=VMPowerState.Value('POWEROFF'))

serialized = req.SerializeToString()
# b'\n\x04name\x10\x02'

deserialized = UpdateVMPowerStateRequest()

deserialized.ParseFromString(serialized)

deserialized
# id: "name"
# status: POWEROFF
```

# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: delete_volume_request.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='delete_volume_request.proto',
  package='',
  syntax='proto3',
  serialized_pb=_b('\n\x1b\x64\x65lete_volume_request.proto\"1\n\x13\x44\x65leteVolumeRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04host\x18\x02 \x01(\tb\x06proto3')
)




_DELETEVOLUMEREQUEST = _descriptor.Descriptor(
  name='DeleteVolumeRequest',
  full_name='DeleteVolumeRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='DeleteVolumeRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='host', full_name='DeleteVolumeRequest.host', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=31,
  serialized_end=80,
)

DESCRIPTOR.message_types_by_name['DeleteVolumeRequest'] = _DELETEVOLUMEREQUEST
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

DeleteVolumeRequest = _reflection.GeneratedProtocolMessageType('DeleteVolumeRequest', (_message.Message,), dict(
  DESCRIPTOR = _DELETEVOLUMEREQUEST,
  __module__ = 'delete_volume_request_pb2'
  # @@protoc_insertion_point(class_scope:DeleteVolumeRequest)
  ))
_sym_db.RegisterMessage(DeleteVolumeRequest)


# @@protoc_insertion_point(module_scope)
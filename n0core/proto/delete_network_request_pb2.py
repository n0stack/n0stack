# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: delete_network_request.proto

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
  name='delete_network_request.proto',
  package='',
  syntax='proto3',
  serialized_pb=_b('\n\x1c\x64\x65lete_network_request.proto\"\"\n\x14\x44\x65leteNetworkRequest\x12\n\n\x02id\x18\x01 \x01(\tb\x06proto3')
)




_DELETENETWORKREQUEST = _descriptor.Descriptor(
  name='DeleteNetworkRequest',
  full_name='DeleteNetworkRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='DeleteNetworkRequest.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
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
  serialized_start=32,
  serialized_end=66,
)

DESCRIPTOR.message_types_by_name['DeleteNetworkRequest'] = _DELETENETWORKREQUEST
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

DeleteNetworkRequest = _reflection.GeneratedProtocolMessageType('DeleteNetworkRequest', (_message.Message,), dict(
  DESCRIPTOR = _DELETENETWORKREQUEST,
  __module__ = 'delete_network_request_pb2'
  # @@protoc_insertion_point(class_scope:DeleteNetworkRequest)
  ))
_sym_db.RegisterMessage(DeleteNetworkRequest)


# @@protoc_insertion_point(module_scope)
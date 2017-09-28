# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: update_network_interface_request.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import vlan_network_setting_pb2 as vlan__network__setting__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='update_network_interface_request.proto',
  package='',
  syntax='proto3',
  serialized_pb=_b('\n&update_network_interface_request.proto\x1a\x1avlan_network_setting.proto\"\x8f\x01\n\x1dUpdateNetworkInterfaceRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x14\n\x0cip_addresses\x18\x02 \x03(\t\x12\x33\n\x14vlan_network_setting\x18\x03 \x01(\x0b\x32\x13.VlanNetworkSettingH\x00\x42\x17\n\x15type_specific_settingb\x06proto3')
  ,
  dependencies=[vlan__network__setting__pb2.DESCRIPTOR,])




_UPDATENETWORKINTERFACEREQUEST = _descriptor.Descriptor(
  name='UpdateNetworkInterfaceRequest',
  full_name='UpdateNetworkInterfaceRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='UpdateNetworkInterfaceRequest.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='ip_addresses', full_name='UpdateNetworkInterfaceRequest.ip_addresses', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='vlan_network_setting', full_name='UpdateNetworkInterfaceRequest.vlan_network_setting', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
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
    _descriptor.OneofDescriptor(
      name='type_specific_setting', full_name='UpdateNetworkInterfaceRequest.type_specific_setting',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=71,
  serialized_end=214,
)

_UPDATENETWORKINTERFACEREQUEST.fields_by_name['vlan_network_setting'].message_type = vlan__network__setting__pb2._VLANNETWORKSETTING
_UPDATENETWORKINTERFACEREQUEST.oneofs_by_name['type_specific_setting'].fields.append(
  _UPDATENETWORKINTERFACEREQUEST.fields_by_name['vlan_network_setting'])
_UPDATENETWORKINTERFACEREQUEST.fields_by_name['vlan_network_setting'].containing_oneof = _UPDATENETWORKINTERFACEREQUEST.oneofs_by_name['type_specific_setting']
DESCRIPTOR.message_types_by_name['UpdateNetworkInterfaceRequest'] = _UPDATENETWORKINTERFACEREQUEST
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

UpdateNetworkInterfaceRequest = _reflection.GeneratedProtocolMessageType('UpdateNetworkInterfaceRequest', (_message.Message,), dict(
  DESCRIPTOR = _UPDATENETWORKINTERFACEREQUEST,
  __module__ = 'update_network_interface_request_pb2'
  # @@protoc_insertion_point(class_scope:UpdateNetworkInterfaceRequest)
  ))
_sym_db.RegisterMessage(UpdateNetworkInterfaceRequest)


# @@protoc_insertion_point(module_scope)
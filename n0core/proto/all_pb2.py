# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: all.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import request_pb2 as request__pb2
try:
  attach__network__interface__request__pb2 = request__pb2.attach__network__interface__request__pb2
except AttributeError:
  attach__network__interface__request__pb2 = request__pb2.attach_network_interface_request_pb2
try:
  clone__vm__request__pb2 = request__pb2.clone__vm__request__pb2
except AttributeError:
  clone__vm__request__pb2 = request__pb2.clone_vm_request_pb2
try:
  create__ipv4__subnet__request__pb2 = request__pb2.create__ipv4__subnet__request__pb2
except AttributeError:
  create__ipv4__subnet__request__pb2 = request__pb2.create_ipv4_subnet_request_pb2
try:
  create__network__request__pb2 = request__pb2.create__network__request__pb2
except AttributeError:
  create__network__request__pb2 = request__pb2.create_network_request_pb2
try:
  create__vm__request__pb2 = request__pb2.create__vm__request__pb2
except AttributeError:
  create__vm__request__pb2 = request__pb2.create_vm_request_pb2
try:
  create__vm__snapshot__request__pb2 = request__pb2.create__vm__snapshot__request__pb2
except AttributeError:
  create__vm__snapshot__request__pb2 = request__pb2.create_vm_snapshot_request_pb2
try:
  create__volume__request__pb2 = request__pb2.create__volume__request__pb2
except AttributeError:
  create__volume__request__pb2 = request__pb2.create_volume_request_pb2
try:
  delete__ipv4__subnet__request__pb2 = request__pb2.delete__ipv4__subnet__request__pb2
except AttributeError:
  delete__ipv4__subnet__request__pb2 = request__pb2.delete_ipv4_subnet_request_pb2
try:
  delete__network__request__pb2 = request__pb2.delete__network__request__pb2
except AttributeError:
  delete__network__request__pb2 = request__pb2.delete_network_request_pb2
try:
  delete__vm__request__pb2 = request__pb2.delete__vm__request__pb2
except AttributeError:
  delete__vm__request__pb2 = request__pb2.delete_vm_request_pb2
try:
  delete__vm__snapshot__request__pb2 = request__pb2.delete__vm__snapshot__request__pb2
except AttributeError:
  delete__vm__snapshot__request__pb2 = request__pb2.delete_vm_snapshot_request_pb2
try:
  delete__volume__request__pb2 = request__pb2.delete__volume__request__pb2
except AttributeError:
  delete__volume__request__pb2 = request__pb2.delete_volume_request_pb2
try:
  detach__network__interface__request__pb2 = request__pb2.detach__network__interface__request__pb2
except AttributeError:
  detach__network__interface__request__pb2 = request__pb2.detach_network_interface_request_pb2
try:
  migrate__vm__request__pb2 = request__pb2.migrate__vm__request__pb2
except AttributeError:
  migrate__vm__request__pb2 = request__pb2.migrate_vm_request_pb2
try:
  update__ipv4__subnet__request__pb2 = request__pb2.update__ipv4__subnet__request__pb2
except AttributeError:
  update__ipv4__subnet__request__pb2 = request__pb2.update_ipv4_subnet_request_pb2
try:
  update__network__interface__request__pb2 = request__pb2.update__network__interface__request__pb2
except AttributeError:
  update__network__interface__request__pb2 = request__pb2.update_network_interface_request_pb2
try:
  update__network__request__pb2 = request__pb2.update__network__request__pb2
except AttributeError:
  update__network__request__pb2 = request__pb2.update_network_request_pb2
try:
  update__vm__power__state__request__pb2 = request__pb2.update__vm__power__state__request__pb2
except AttributeError:
  update__vm__power__state__request__pb2 = request__pb2.update_vm_power_state_request_pb2
try:
  update__vm__request__pb2 = request__pb2.update__vm__request__pb2
except AttributeError:
  update__vm__request__pb2 = request__pb2.update_vm_request_pb2
try:
  network__type__pb2 = request__pb2.network__type__pb2
except AttributeError:
  network__type__pb2 = request__pb2.network_type_pb2
try:
  vm__power__state__pb2 = request__pb2.vm__power__state__pb2
except AttributeError:
  vm__power__state__pb2 = request__pb2.vm_power_state_pb2
try:
  vlan__network__setting__pb2 = request__pb2.vlan__network__setting__pb2
except AttributeError:
  vlan__network__setting__pb2 = request__pb2.vlan_network_setting_pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='all.proto',
  package='',
  syntax='proto3',
  serialized_pb=_b('\n\tall.proto\x1a\rrequest.protob\x06proto3')
  ,
  dependencies=[request__pb2.DESCRIPTOR,])



_sym_db.RegisterFileDescriptor(DESCRIPTOR)


# @@protoc_insertion_point(module_scope)
# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [n0stack/infra/v1alpha/block_storage.proto](#n0stack/infra/v1alpha/block_storage.proto)
    - [BlockStorage](#n0stack.infra.v1alpha.BlockStorage)
    - [BlockStorage.AnnotationsEntry](#n0stack.infra.v1alpha.BlockStorage.AnnotationsEntry)
    - [BlockStorage.LabelsEntry](#n0stack.infra.v1alpha.BlockStorage.LabelsEntry)
    - [CancelBlockStorageOperationRequest](#n0stack.infra.v1alpha.CancelBlockStorageOperationRequest)
    - [CopyBlockStorageRequest](#n0stack.infra.v1alpha.CopyBlockStorageRequest)
    - [CreateBlockStorageRequest](#n0stack.infra.v1alpha.CreateBlockStorageRequest)
    - [DeleteBlockStorageRequest](#n0stack.infra.v1alpha.DeleteBlockStorageRequest)
    - [FetchBlockStorageRequest](#n0stack.infra.v1alpha.FetchBlockStorageRequest)
    - [GetBlockStorageRequest](#n0stack.infra.v1alpha.GetBlockStorageRequest)
    - [ListBlockStoragesRequest](#n0stack.infra.v1alpha.ListBlockStoragesRequest)
    - [ListBlockStoragesResponse](#n0stack.infra.v1alpha.ListBlockStoragesResponse)
    - [ReleaseBlockStorageRequest](#n0stack.infra.v1alpha.ReleaseBlockStorageRequest)
    - [UpdateBlockStorageRequest](#n0stack.infra.v1alpha.UpdateBlockStorageRequest)
    - [UploadBlockStorageRequest](#n0stack.infra.v1alpha.UploadBlockStorageRequest)
    - [UploadBlockStorageResponse](#n0stack.infra.v1alpha.UploadBlockStorageResponse)
    - [UseBlockStorageRequest](#n0stack.infra.v1alpha.UseBlockStorageRequest)
  
    - [BlockStorage.BlockStorageState](#n0stack.infra.v1alpha.BlockStorage.BlockStorageState)
  
  
    - [BlockStorageService](#n0stack.infra.v1alpha.BlockStorageService)
  

- [n0stack/infra/v1alpha/network.proto](#n0stack/infra/v1alpha/network.proto)
    - [CancelNetworkOperationRequest](#n0stack.infra.v1alpha.CancelNetworkOperationRequest)
    - [CreateNetworkRequest](#n0stack.infra.v1alpha.CreateNetworkRequest)
    - [DeleteNetworkRequest](#n0stack.infra.v1alpha.DeleteNetworkRequest)
    - [GetNetworkRequest](#n0stack.infra.v1alpha.GetNetworkRequest)
    - [ListNetworksRequest](#n0stack.infra.v1alpha.ListNetworksRequest)
    - [ListNetworksResponse](#n0stack.infra.v1alpha.ListNetworksResponse)
    - [Network](#n0stack.infra.v1alpha.Network)
    - [Network.AnnotationsEntry](#n0stack.infra.v1alpha.Network.AnnotationsEntry)
    - [Network.LabelsEntry](#n0stack.infra.v1alpha.Network.LabelsEntry)
    - [NetworkInterface](#n0stack.infra.v1alpha.NetworkInterface)
    - [NetworkInterface.AnnotationsEntry](#n0stack.infra.v1alpha.NetworkInterface.AnnotationsEntry)
    - [NetworkInterface.LabelsEntry](#n0stack.infra.v1alpha.NetworkInterface.LabelsEntry)
    - [ReleaseNetworkInterfaceRequest](#n0stack.infra.v1alpha.ReleaseNetworkInterfaceRequest)
    - [ReserveNetworkInterfaceRequest](#n0stack.infra.v1alpha.ReserveNetworkInterfaceRequest)
    - [UpdateNetworkRequest](#n0stack.infra.v1alpha.UpdateNetworkRequest)
  
    - [Network.NetworkState](#n0stack.infra.v1alpha.Network.NetworkState)
  
  
    - [NetworkService](#n0stack.infra.v1alpha.NetworkService)
  

- [n0stack/infra/v1alpha/virtual_machine.proto](#n0stack/infra/v1alpha/virtual_machine.proto)
    - [BootVirtualMachineRequest](#n0stack.infra.v1alpha.BootVirtualMachineRequest)
    - [CancelVirtualMachineOperationRequest](#n0stack.infra.v1alpha.CancelVirtualMachineOperationRequest)
    - [ChangeVirtualMachineRunningStateRequest](#n0stack.infra.v1alpha.ChangeVirtualMachineRunningStateRequest)
    - [CreateVirtualMachineRequest](#n0stack.infra.v1alpha.CreateVirtualMachineRequest)
    - [DeleteVirtualMachineRequest](#n0stack.infra.v1alpha.DeleteVirtualMachineRequest)
    - [GetVirtualMachineRequest](#n0stack.infra.v1alpha.GetVirtualMachineRequest)
    - [ListVirtualMachinesRequest](#n0stack.infra.v1alpha.ListVirtualMachinesRequest)
    - [ListVirtualMachinesResponse](#n0stack.infra.v1alpha.ListVirtualMachinesResponse)
    - [OpenConsoleRequest](#n0stack.infra.v1alpha.OpenConsoleRequest)
    - [OpenConsoleResponse](#n0stack.infra.v1alpha.OpenConsoleResponse)
    - [RebootVirtualMachineRequest](#n0stack.infra.v1alpha.RebootVirtualMachineRequest)
    - [ShutdownVirtualMachineRequest](#n0stack.infra.v1alpha.ShutdownVirtualMachineRequest)
    - [UpdateVirtualMachineRequest](#n0stack.infra.v1alpha.UpdateVirtualMachineRequest)
    - [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine)
    - [VirtualMachine.AnnotationsEntry](#n0stack.infra.v1alpha.VirtualMachine.AnnotationsEntry)
    - [VirtualMachine.LabelsEntry](#n0stack.infra.v1alpha.VirtualMachine.LabelsEntry)
    - [VirtualMachine.VirtualMachineNIC](#n0stack.infra.v1alpha.VirtualMachine.VirtualMachineNIC)
  
    - [VirtualMachine.VirtualMachineState](#n0stack.infra.v1alpha.VirtualMachine.VirtualMachineState)
  
  
    - [VirtualMachineService](#n0stack.infra.v1alpha.VirtualMachineService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="n0stack/infra/v1alpha/block_storage.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/infra/v1alpha/block_storage.proto



<a name="n0stack.infra.v1alpha.BlockStorage"></a>

### BlockStorage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| project | [string](#string) |  |  |
| annotations | [BlockStorage.AnnotationsEntry](#n0stack.infra.v1alpha.BlockStorage.AnnotationsEntry) | repeated | Annotations can store metadata used by the system for control. In particular, implementation-dependent fields that can not be set as protobuf fields are targeted. The control specified by n0stack may delete metadata specified by the user. |
| labels | [BlockStorage.LabelsEntry](#n0stack.infra.v1alpha.BlockStorage.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| operation | [n0stack.protobuf.Operation](#n0stack.protobuf.Operation) |  |  |
| state | [BlockStorage.BlockStorageState](#n0stack.infra.v1alpha.BlockStorage.BlockStorageState) |  |  |
| bytes | [uint64](#uint64) |  |  |
| is_cd | [bool](#bool) |  |  |
| download_url | [string](#string) |  |  |
| in_use | [bool](#bool) |  |  |
| in_use_reason | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.BlockStorage.AnnotationsEntry"></a>

### BlockStorage.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.BlockStorage.LabelsEntry"></a>

### BlockStorage.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.CancelBlockStorageOperationRequest"></a>

### CancelBlockStorageOperationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.CopyBlockStorageRequest"></a>

### CopyBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |
| source_block_storage_name | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.CreateBlockStorageRequest"></a>

### CreateBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |






<a name="n0stack.infra.v1alpha.DeleteBlockStorageRequest"></a>

### DeleteBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.FetchBlockStorageRequest"></a>

### FetchBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |
| source_url | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.GetBlockStorageRequest"></a>

### GetBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListBlockStoragesRequest"></a>

### ListBlockStoragesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListBlockStoragesResponse"></a>

### ListBlockStoragesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storages | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) | repeated |  |






<a name="n0stack.infra.v1alpha.ReleaseBlockStorageRequest"></a>

### ReleaseBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.UpdateBlockStorageRequest"></a>

### UpdateBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |






<a name="n0stack.infra.v1alpha.UploadBlockStorageRequest"></a>

### UploadBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |






<a name="n0stack.infra.v1alpha.UploadBlockStorageResponse"></a>

### UploadBlockStorageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| block_storage | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |  |
| upload_url | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.UseBlockStorageRequest"></a>

### UseBlockStorageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| in_use_reason | [string](#string) |  |  |





 


<a name="n0stack.infra.v1alpha.BlockStorage.BlockStorageState"></a>

### BlockStorage.BlockStorageState


| Name | Number | Description |
| ---- | ------ | ----------- |
| BLOCK_STORAGE_UNSPECIFIED | 0 |  |
| AVAILABLE | 1 |  |
| DELETED | 2 |  |
| CREATING | 16 | standard unsteady state |
| DELETING | 17 |  |
| FETCHING | 32 |  |
| CLONING | 33 |  |
| UPLOADING | 34 |  |


 

 


<a name="n0stack.infra.v1alpha.BlockStorageService"></a>

### BlockStorageService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateBlockStorage | [CreateBlockStorageRequest](#n0stack.infra.v1alpha.CreateBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| CloneBlockStorage | [CopyBlockStorageRequest](#n0stack.infra.v1alpha.CopyBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| FetchBlockStorage | [FetchBlockStorageRequest](#n0stack.infra.v1alpha.FetchBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| UploadBlockStorage | [UploadBlockStorageRequest](#n0stack.infra.v1alpha.UploadBlockStorageRequest) | [UploadBlockStorageResponse](#n0stack.infra.v1alpha.UploadBlockStorageResponse) |  |
| ListBlockStorages | [ListBlockStoragesRequest](#n0stack.infra.v1alpha.ListBlockStoragesRequest) | [ListBlockStoragesResponse](#n0stack.infra.v1alpha.ListBlockStoragesResponse) |  |
| GetBlockStorage | [GetBlockStorageRequest](#n0stack.infra.v1alpha.GetBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| UpdateBlockStorage | [UpdateBlockStorageRequest](#n0stack.infra.v1alpha.UpdateBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| DeleteBlockStorage | [DeleteBlockStorageRequest](#n0stack.infra.v1alpha.DeleteBlockStorageRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| CancelBlockStorageOperation | [CancelBlockStorageOperationRequest](#n0stack.infra.v1alpha.CancelBlockStorageOperationRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| ProposeBlockStorageOperation | [.n0stack.protobuf.ProposeOperationRequest](#n0stack.protobuf.ProposeOperationRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| UseBlockStorage | [UseBlockStorageRequest](#n0stack.infra.v1alpha.UseBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |
| ReleaseBlockStorage | [ReleaseBlockStorageRequest](#n0stack.infra.v1alpha.ReleaseBlockStorageRequest) | [BlockStorage](#n0stack.infra.v1alpha.BlockStorage) |  |

 



<a name="n0stack/infra/v1alpha/network.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/infra/v1alpha/network.proto



<a name="n0stack.infra.v1alpha.CancelNetworkOperationRequest"></a>

### CancelNetworkOperationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.CreateNetworkRequest"></a>

### CreateNetworkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [Network](#n0stack.infra.v1alpha.Network) |  |  |






<a name="n0stack.infra.v1alpha.DeleteNetworkRequest"></a>

### DeleteNetworkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.GetNetworkRequest"></a>

### GetNetworkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListNetworksRequest"></a>

### ListNetworksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListNetworksResponse"></a>

### ListNetworksResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| networks | [Network](#n0stack.infra.v1alpha.Network) | repeated |  |






<a name="n0stack.infra.v1alpha.Network"></a>

### Network



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| project | [string](#string) |  |  |
| annotations | [Network.AnnotationsEntry](#n0stack.infra.v1alpha.Network.AnnotationsEntry) | repeated | The annotations parameter can store metadata used by the n0stack system. The n0stack operations may modify this defined by the user without any notice. In particular, this targets domain specific parameters, which cannot be used by many users. |
| labels | [Network.LabelsEntry](#n0stack.infra.v1alpha.Network.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| state | [Network.NetworkState](#n0stack.infra.v1alpha.Network.NetworkState) |  |  |
| ipv4_cidr | [string](#string) |  |  |
| ipv6_cidr | [string](#string) |  |  |
| network_interfaces | [NetworkInterface](#n0stack.infra.v1alpha.NetworkInterface) | repeated |  |






<a name="n0stack.infra.v1alpha.Network.AnnotationsEntry"></a>

### Network.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.Network.LabelsEntry"></a>

### Network.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.NetworkInterface"></a>

### NetworkInterface



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| annotations | [NetworkInterface.AnnotationsEntry](#n0stack.infra.v1alpha.NetworkInterface.AnnotationsEntry) | repeated |  |
| labels | [NetworkInterface.LabelsEntry](#n0stack.infra.v1alpha.NetworkInterface.LabelsEntry) | repeated |  |
| hardware_address | [string](#string) |  | Network の中でユニークであること |
| ipv4_address | [string](#string) |  | Network の CIDR に含まれていること |
| ipv6_address | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.NetworkInterface.AnnotationsEntry"></a>

### NetworkInterface.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.NetworkInterface.LabelsEntry"></a>

### NetworkInterface.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ReleaseNetworkInterfaceRequest"></a>

### ReleaseNetworkInterfaceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network_name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| network_interface_name | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ReserveNetworkInterfaceRequest"></a>

### ReserveNetworkInterfaceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network_name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| network_interface | [NetworkInterface](#n0stack.infra.v1alpha.NetworkInterface) |  |  |






<a name="n0stack.infra.v1alpha.UpdateNetworkRequest"></a>

### UpdateNetworkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [Network](#n0stack.infra.v1alpha.Network) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |





 


<a name="n0stack.infra.v1alpha.Network.NetworkState"></a>

### Network.NetworkState


| Name | Number | Description |
| ---- | ------ | ----------- |
| NETWORK_UNSPECIFIED | 0 | falied state because failed some process on API. |
| AVAILABLE | 1 | steady state |
| CREATING | 16 | standard unsteady state |
| DELETING | 17 |  |
| RESERVING_NETWORK_INTERFACE | 32 |  |
| RELEASING_NETWORK_INTERFACE | 33 |  |


 

 


<a name="n0stack.infra.v1alpha.NetworkService"></a>

### NetworkService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateNetwork | [CreateNetworkRequest](#n0stack.infra.v1alpha.CreateNetworkRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| ListNetworks | [ListNetworksRequest](#n0stack.infra.v1alpha.ListNetworksRequest) | [ListNetworksResponse](#n0stack.infra.v1alpha.ListNetworksResponse) |  |
| GetNetwork | [GetNetworkRequest](#n0stack.infra.v1alpha.GetNetworkRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| UpdateNetwork | [UpdateNetworkRequest](#n0stack.infra.v1alpha.UpdateNetworkRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| DeleteNetwork | [DeleteNetworkRequest](#n0stack.infra.v1alpha.DeleteNetworkRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| CancelNetworkOperation | [CancelNetworkOperationRequest](#n0stack.infra.v1alpha.CancelNetworkOperationRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| ProposeNetworkOperation | [.n0stack.protobuf.ProposeOperationRequest](#n0stack.protobuf.ProposeOperationRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| ReserveNetworkInterface | [ReserveNetworkInterfaceRequest](#n0stack.infra.v1alpha.ReserveNetworkInterfaceRequest) | [Network](#n0stack.infra.v1alpha.Network) |  |
| ReleaseNetworkInterface | [ReleaseNetworkInterfaceRequest](#n0stack.infra.v1alpha.ReleaseNetworkInterfaceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



<a name="n0stack/infra/v1alpha/virtual_machine.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/infra/v1alpha/virtual_machine.proto



<a name="n0stack.infra.v1alpha.BootVirtualMachineRequest"></a>

### BootVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.CancelVirtualMachineOperationRequest"></a>

### CancelVirtualMachineOperationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ChangeVirtualMachineRunningStateRequest"></a>

### ChangeVirtualMachineRunningStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| is_running | [bool](#bool) |  |  |






<a name="n0stack.infra.v1alpha.CreateVirtualMachineRequest"></a>

### CreateVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| virtual_machine | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) |  |  |






<a name="n0stack.infra.v1alpha.DeleteVirtualMachineRequest"></a>

### DeleteVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.GetVirtualMachineRequest"></a>

### GetVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListVirtualMachinesRequest"></a>

### ListVirtualMachinesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.ListVirtualMachinesResponse"></a>

### ListVirtualMachinesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| virtual_machines | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | repeated |  |






<a name="n0stack.infra.v1alpha.OpenConsoleRequest"></a>

### OpenConsoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.OpenConsoleResponse"></a>

### OpenConsoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| console_url | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.RebootVirtualMachineRequest"></a>

### RebootVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| hard | [bool](#bool) |  |  |






<a name="n0stack.infra.v1alpha.ShutdownVirtualMachineRequest"></a>

### ShutdownVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| hard | [bool](#bool) |  |  |






<a name="n0stack.infra.v1alpha.UpdateVirtualMachineRequest"></a>

### UpdateVirtualMachineRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| virtual_machine | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |






<a name="n0stack.infra.v1alpha.VirtualMachine"></a>

### VirtualMachine



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| project | [string](#string) |  |  |
| annotations | [VirtualMachine.AnnotationsEntry](#n0stack.infra.v1alpha.VirtualMachine.AnnotationsEntry) | repeated | The annotations parameter can store metadata used by the n0stack system. The n0stack operations may modify this defined by the user without any notice. In particular, this targets domain specific parameters, which cannot be used by many users. |
| labels | [VirtualMachine.LabelsEntry](#n0stack.infra.v1alpha.VirtualMachine.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| operation | [n0stack.protobuf.Operation](#n0stack.protobuf.Operation) |  |  |
| state | [VirtualMachine.VirtualMachineState](#n0stack.infra.v1alpha.VirtualMachine.VirtualMachineState) |  |  |
| uuid | [string](#string) |  |  |
| cpu_core | [uint32](#uint32) |  |  |
| memory_bytes | [uint64](#uint64) |  |  |
| block_storage_names | [string](#string) | repeated |  |
| nics | [VirtualMachine.VirtualMachineNIC](#n0stack.infra.v1alpha.VirtualMachine.VirtualMachineNIC) | repeated |  |
| is_running | [bool](#bool) |  |  |
| login_user_names | [string](#string) | repeated |  |
| login_group_names | [string](#string) | repeated |  |






<a name="n0stack.infra.v1alpha.VirtualMachine.AnnotationsEntry"></a>

### VirtualMachine.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.VirtualMachine.LabelsEntry"></a>

### VirtualMachine.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.infra.v1alpha.VirtualMachine.VirtualMachineNIC"></a>

### VirtualMachine.VirtualMachineNIC



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network_name | [string](#string) |  |  |
| network_interface_name | [string](#string) |  |  |





 


<a name="n0stack.infra.v1alpha.VirtualMachine.VirtualMachineState"></a>

### VirtualMachine.VirtualMachineState


| Name | Number | Description |
| ---- | ------ | ----------- |
| VIRTUAL_MACHINE_UNSPECIFIED | 0 |  |
| AVAILABLE | 1 | steady state |
| DELETED | 2 |  |
| CREATING | 16 | standard unsteady state |
| DELETING | 17 |  |
| BOOTING | 32 |  |
| REBOOTING | 33 |  |
| STOPPING | 34 |  |


 

 


<a name="n0stack.infra.v1alpha.VirtualMachineService"></a>

### VirtualMachineService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateVirtualMachine | [CreateVirtualMachineRequest](#n0stack.infra.v1alpha.CreateVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| ListVirtualMachines | [ListVirtualMachinesRequest](#n0stack.infra.v1alpha.ListVirtualMachinesRequest) | [ListVirtualMachinesResponse](#n0stack.infra.v1alpha.ListVirtualMachinesResponse) | summary: description: required_roles: - project: n0stack role: virtual_machine:read - project: * role: n0stack:virtual_machine:read errors: NotFound: Return if the virtual machine resource is not found |
| GetVirtualMachine | [GetVirtualMachineRequest](#n0stack.infra.v1alpha.GetVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:read - project: * role: n0stack:virtual_machine:read errors: NotFound: Return if the virtual machine resource is not found |
| UpdateVirtualMachine | [UpdateVirtualMachineRequest](#n0stack.infra.v1alpha.UpdateVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) |  |
| DeleteVirtualMachine | [DeleteVirtualMachineRequest](#n0stack.infra.v1alpha.DeleteVirtualMachineRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| BootVirtualMachine | [BootVirtualMachineRequest](#n0stack.infra.v1alpha.BootVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| RebootVirtualMachine | [RebootVirtualMachineRequest](#n0stack.infra.v1alpha.RebootVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| ShutdownVirtualMachine | [ShutdownVirtualMachineRequest](#n0stack.infra.v1alpha.ShutdownVirtualMachineRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| CancelVirtualMachineOperation | [CancelVirtualMachineOperationRequest](#n0stack.infra.v1alpha.CancelVirtualMachineOperationRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: description: required_roles: - project: n0stack role: virtual_machine:write - project: * role: n0stack:virtual_machine:write errors: NotFound: Return if the virtual machine resource is not found |
| OpenVirtualMachineConsole | [OpenConsoleRequest](#n0stack.infra.v1alpha.OpenConsoleRequest) | [OpenConsoleResponse](#n0stack.infra.v1alpha.OpenConsoleResponse) | summary: OpenVirtualMachineConsole returns a URL to open the console of the virtual machine. description: required_roles: - project: n0stack role: virtual_machine:read - project: * role: n0stack:virtual_machine:read errors: NotFound: Return if the virtual machine resource is not found |
| ProposeVirtualMachineOperation | [.n0stack.protobuf.ProposeOperationRequest](#n0stack.protobuf.ProposeOperationRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: Propose*Operation propose the resource operation description:| You can propose the resource operation if the resource state is unsteady. The API ensures that only one operator is operating using lock with operating_for_seconds. If the operation is failed, set failed backoff to operating_for_seconds. required_roles: - project: n0stack role: virtual_machine:operate errors: NotFound: Return if the virtual machine resource is not found |
| ChangeVirtualMachineRunningState | [ChangeVirtualMachineRunningStateRequest](#n0stack.infra.v1alpha.ChangeVirtualMachineRunningStateRequest) | [VirtualMachine](#n0stack.infra.v1alpha.VirtualMachine) | summary: ChangeVirtualMachineRunningState change the is_running parameter of the virtual machine description:| Set the is_running parameter to true if the virtual machine is running. required_roles: - project: n0stack role: virtual_machine:operate errors: NotFound: Return if the virtual machine resource is not found |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |


# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [n0stack/protobuf/operation.proto](#n0stack/protobuf/operation.proto)
    - [Operation](#n0stack.protobuf.Operation)
    - [Operation.Log](#n0stack.protobuf.Operation.Log)
    - [ProposeOperationRequest](#n0stack.protobuf.ProposeOperationRequest)
  
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="n0stack/protobuf/operation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/protobuf/operation.proto



<a name="n0stack.protobuf.Operation"></a>

### Operation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation_request | [google.protobuf.Any](#google.protobuf.Any) |  | 例えばShutdownVirtualMachineが実行され、STOPPINGステータスになった場合ShutdownVirtualMachineRequestが格納される |
| locked_until | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | locked_untilにセットされた時間まで他のoperatorがリソースを操作してはいけない |
| is_disabled_cancel | [bool](#bool) |  | is_disabled_cancel == true &amp;&amp; !locked_until.is_expired() の間は Cancel*() ができない is_disabled_cancel == true &amp;&amp; locked_until.is_expired() の場合、operatorが死んでいる |
| proposing_peer_token | [string](#string) |  | ミスオペを防止するためのゆるくてもチェック機構が必要になるので、peerがユニークなトークンを提案し、ロックされている間はそれを検証する |
| logs | [Operation.Log](#n0stack.protobuf.Operation.Log) | repeated |  |






<a name="n0stack.protobuf.Operation.Log"></a>

### Operation.Log



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ts | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| message | [string](#string) |  |  |






<a name="n0stack.protobuf.ProposeOperationRequest"></a>

### ProposeOperationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |
| lock_for_seconds | [uint32](#uint32) |  |  |
| log_message | [string](#string) |  |  |
| is_disabled_cancel | [bool](#bool) |  |  |
| proposing_peer_token | [string](#string) |  | プロセス開始時にハッシュ値などを自動生成し、それを指定すること |





 

 

 

 



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


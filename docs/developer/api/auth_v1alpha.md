# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [auth/v1alpha/authentication.proto](#auth/v1alpha/authentication.proto)
    - [GetAuthenticationTokenPublicKeyRequest](#n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyRequest)
    - [GetAuthenticationTokenPublicKeyResponse](#n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyResponse)
    - [PublicKeyAuthenticateRequest](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest)
    - [PublicKeyAuthenticateRequest.Response](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Response)
    - [PublicKeyAuthenticateRequest.Start](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Start)
    - [PublicKeyAuthenticateResponse](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse)
    - [PublicKeyAuthenticateResponse.Challenge](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Challenge)
    - [PublicKeyAuthenticateResponse.Result](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Result)
  
  
  
    - [AuthenticationService](#n0stack.auth.v1alpha.AuthenticationService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="auth/v1alpha/authentication.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auth/v1alpha/authentication.proto



<a name="n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyRequest"></a>

### GetAuthenticationTokenPublicKeyRequest







<a name="n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyResponse"></a>

### GetAuthenticationTokenPublicKeyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| public_key | [string](#string) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateRequest"></a>

### PublicKeyAuthenticateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [PublicKeyAuthenticateRequest.Start](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Start) |  |  |
| response | [PublicKeyAuthenticateRequest.Response](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Response) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Response"></a>

### PublicKeyAuthenticateRequest.Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| challenge_token | [string](#string) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateRequest.Start"></a>

### PublicKeyAuthenticateRequest.Start



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_name | [string](#string) |  |  |
| public_key | [string](#string) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateResponse"></a>

### PublicKeyAuthenticateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| challenge | [PublicKeyAuthenticateResponse.Challenge](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Challenge) |  |  |
| result | [PublicKeyAuthenticateResponse.Result](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Result) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Challenge"></a>

### PublicKeyAuthenticateResponse.Challenge



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| challenge | [bytes](#bytes) |  |  |






<a name="n0stack.auth.v1alpha.PublicKeyAuthenticateResponse.Result"></a>

### PublicKeyAuthenticateResponse.Result



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authentication_token | [string](#string) |  |  |





 

 

 


<a name="n0stack.auth.v1alpha.AuthenticationService"></a>

### AuthenticationService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAuthenticationTokenPublicKey | [GetAuthenticationTokenPublicKeyRequest](#n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyRequest) | [GetAuthenticationTokenPublicKeyResponse](#n0stack.auth.v1alpha.GetAuthenticationTokenPublicKeyResponse) |  |
| PublicKeyAuthenticate | [PublicKeyAuthenticateRequest](#n0stack.auth.v1alpha.PublicKeyAuthenticateRequest) stream | [PublicKeyAuthenticateResponse](#n0stack.auth.v1alpha.PublicKeyAuthenticateResponse) stream | authentication method for grpc |

 



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


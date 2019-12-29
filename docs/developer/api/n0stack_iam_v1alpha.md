# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [n0stack/iam/v1alpha/project.proto](#n0stack/iam/v1alpha/project.proto)
    - [AddProjectMembershipRequest](#n0stack.iam.v1alpha.AddProjectMembershipRequest)
    - [CreateProjectRequest](#n0stack.iam.v1alpha.CreateProjectRequest)
    - [DeleteProjectMembershipRequest](#n0stack.iam.v1alpha.DeleteProjectMembershipRequest)
    - [DeleteProjectRequest](#n0stack.iam.v1alpha.DeleteProjectRequest)
    - [GetProjectRequest](#n0stack.iam.v1alpha.GetProjectRequest)
    - [ListProjectsRequest](#n0stack.iam.v1alpha.ListProjectsRequest)
    - [ListProjectsResponse](#n0stack.iam.v1alpha.ListProjectsResponse)
    - [Project](#n0stack.iam.v1alpha.Project)
    - [Project.AnnotationsEntry](#n0stack.iam.v1alpha.Project.AnnotationsEntry)
    - [Project.LabelsEntry](#n0stack.iam.v1alpha.Project.LabelsEntry)
    - [Project.MembershipEntry](#n0stack.iam.v1alpha.Project.MembershipEntry)
    - [UpdateProjectRequest](#n0stack.iam.v1alpha.UpdateProjectRequest)
  
    - [ProjectMembership](#n0stack.iam.v1alpha.ProjectMembership)
  
  
    - [ProjectService](#n0stack.iam.v1alpha.ProjectService)
  

- [n0stack/iam/v1alpha/service_account.proto](#n0stack/iam/v1alpha/service_account.proto)
    - [CreateServiceAccountRequest](#n0stack.iam.v1alpha.CreateServiceAccountRequest)
    - [DeleteServiceAccountRequest](#n0stack.iam.v1alpha.DeleteServiceAccountRequest)
    - [GetServiceAccountRequest](#n0stack.iam.v1alpha.GetServiceAccountRequest)
    - [ListServiceAccountRequest](#n0stack.iam.v1alpha.ListServiceAccountRequest)
    - [ListServiceAccountResponse](#n0stack.iam.v1alpha.ListServiceAccountResponse)
    - [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount)
    - [ServiceAccount.AnnotationsEntry](#n0stack.iam.v1alpha.ServiceAccount.AnnotationsEntry)
    - [ServiceAccount.LabelsEntry](#n0stack.iam.v1alpha.ServiceAccount.LabelsEntry)
    - [ServiceAccount.PublicKeysEntry](#n0stack.iam.v1alpha.ServiceAccount.PublicKeysEntry)
    - [UpdateServiceAccountRequest](#n0stack.iam.v1alpha.UpdateServiceAccountRequest)
  
  
  
    - [ServiceAccountService](#n0stack.iam.v1alpha.ServiceAccountService)
  

- [n0stack/iam/v1alpha/user.proto](#n0stack/iam/v1alpha/user.proto)
    - [CreateUserRequest](#n0stack.iam.v1alpha.CreateUserRequest)
    - [DeleteUserRequest](#n0stack.iam.v1alpha.DeleteUserRequest)
    - [GetUserRequest](#n0stack.iam.v1alpha.GetUserRequest)
    - [UpdateUserRequest](#n0stack.iam.v1alpha.UpdateUserRequest)
    - [User](#n0stack.iam.v1alpha.User)
    - [User.AnnotationsEntry](#n0stack.iam.v1alpha.User.AnnotationsEntry)
    - [User.LabelsEntry](#n0stack.iam.v1alpha.User.LabelsEntry)
    - [User.PublicKeysEntry](#n0stack.iam.v1alpha.User.PublicKeysEntry)
  
  
  
    - [UserService](#n0stack.iam.v1alpha.UserService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="n0stack/iam/v1alpha/project.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/iam/v1alpha/project.proto



<a name="n0stack.iam.v1alpha.AddProjectMembershipRequest"></a>

### AddProjectMembershipRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_name | [string](#string) |  |  |
| user_name | [string](#string) |  |  |
| membership | [ProjectMembership](#n0stack.iam.v1alpha.ProjectMembership) |  |  |






<a name="n0stack.iam.v1alpha.CreateProjectRequest"></a>

### CreateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#n0stack.iam.v1alpha.Project) |  |  |






<a name="n0stack.iam.v1alpha.DeleteProjectMembershipRequest"></a>

### DeleteProjectMembershipRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_name | [string](#string) |  |  |
| user_name | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.DeleteProjectRequest"></a>

### DeleteProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.GetProjectRequest"></a>

### GetProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.ListProjectsRequest"></a>

### ListProjectsRequest







<a name="n0stack.iam.v1alpha.ListProjectsResponse"></a>

### ListProjectsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| projects | [Project](#n0stack.iam.v1alpha.Project) | repeated |  |






<a name="n0stack.iam.v1alpha.Project"></a>

### Project



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| annotations | [Project.AnnotationsEntry](#n0stack.iam.v1alpha.Project.AnnotationsEntry) | repeated | Annotations can store metadata used by the system for control. In particular, implementation-dependent fields that can not be set as protobuf fields are targeted. The control specified by n0stack may delete metadata specified by the user. |
| labels | [Project.LabelsEntry](#n0stack.iam.v1alpha.Project.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| display_name | [string](#string) |  |  |
| membership | [Project.MembershipEntry](#n0stack.iam.v1alpha.Project.MembershipEntry) | repeated |  |






<a name="n0stack.iam.v1alpha.Project.AnnotationsEntry"></a>

### Project.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.Project.LabelsEntry"></a>

### Project.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.Project.MembershipEntry"></a>

### Project.MembershipEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [ProjectMembership](#n0stack.iam.v1alpha.ProjectMembership) |  |  |






<a name="n0stack.iam.v1alpha.UpdateProjectRequest"></a>

### UpdateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#n0stack.iam.v1alpha.Project) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |





 


<a name="n0stack.iam.v1alpha.ProjectMembership"></a>

### ProjectMembership


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROJECT_MEMBERSHIP_UNSPECIFIED | 0 |  |
| OWNER | 1 | Owners have all of permissions. |
| MEMBER | 2 | Members have only assined permissions by Roles. |


 

 


<a name="n0stack.iam.v1alpha.ProjectService"></a>

### ProjectService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListProjects | [ListProjectsRequest](#n0stack.iam.v1alpha.ListProjectsRequest) | [ListProjectsResponse](#n0stack.iam.v1alpha.ListProjectsResponse) | あるユーザーがメンバーであるプロジェクトすべてを返す errors: NotFound: どのプロジェクトにも属していない場合 |
| GetProject | [GetProjectRequest](#n0stack.iam.v1alpha.GetProjectRequest) | [Project](#n0stack.iam.v1alpha.Project) | Summary: プロジェクトの詳細を取得する errors: NotFound: memberじゃない場合 Unauthorized: ログインしていない場合 |
| CreateProject | [CreateProjectRequest](#n0stack.iam.v1alpha.CreateProjectRequest) | [Project](#n0stack.iam.v1alpha.Project) | ログインしているユーザーがオーナーとなるプロジェクトを作成する errors: Unauthorized: ログインしていない場合 |
| UpdateProject | [UpdateProjectRequest](#n0stack.iam.v1alpha.UpdateProjectRequest) | [Project](#n0stack.iam.v1alpha.Project) |  |
| DeleteProject | [DeleteProjectRequest](#n0stack.iam.v1alpha.DeleteProjectRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AddProjectMembership | [AddProjectMembershipRequest](#n0stack.iam.v1alpha.AddProjectMembershipRequest) | [Project](#n0stack.iam.v1alpha.Project) |  |
| DeleteProjectMembership | [DeleteProjectMembershipRequest](#n0stack.iam.v1alpha.DeleteProjectMembershipRequest) | [Project](#n0stack.iam.v1alpha.Project) |  |

 



<a name="n0stack/iam/v1alpha/service_account.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/iam/v1alpha/service_account.proto



<a name="n0stack.iam.v1alpha.CreateServiceAccountRequest"></a>

### CreateServiceAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_account | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) |  |  |






<a name="n0stack.iam.v1alpha.DeleteServiceAccountRequest"></a>

### DeleteServiceAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.GetServiceAccountRequest"></a>

### GetServiceAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.ListServiceAccountRequest"></a>

### ListServiceAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.ListServiceAccountResponse"></a>

### ListServiceAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_accounts | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) | repeated |  |






<a name="n0stack.iam.v1alpha.ServiceAccount"></a>

### ServiceAccount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| project | [string](#string) |  |  |
| annotations | [ServiceAccount.AnnotationsEntry](#n0stack.iam.v1alpha.ServiceAccount.AnnotationsEntry) | repeated | Annotations can store metadata used by the system for control. In particular, implementation-dependent fields that can not be set as protobuf fields are targeted. The control specified by n0stack may delete metadata specified by the user. |
| labels | [ServiceAccount.LabelsEntry](#n0stack.iam.v1alpha.ServiceAccount.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| public_keys | [ServiceAccount.PublicKeysEntry](#n0stack.iam.v1alpha.ServiceAccount.PublicKeysEntry) | repeated |  |






<a name="n0stack.iam.v1alpha.ServiceAccount.AnnotationsEntry"></a>

### ServiceAccount.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.ServiceAccount.LabelsEntry"></a>

### ServiceAccount.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.ServiceAccount.PublicKeysEntry"></a>

### ServiceAccount.PublicKeysEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.UpdateServiceAccountRequest"></a>

### UpdateServiceAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_account | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |





 

 

 


<a name="n0stack.iam.v1alpha.ServiceAccountService"></a>

### ServiceAccountService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateServiceAccount | [CreateServiceAccountRequest](#n0stack.iam.v1alpha.CreateServiceAccountRequest) | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) |  |
| ListServiceAccount | [ListServiceAccountRequest](#n0stack.iam.v1alpha.ListServiceAccountRequest) | [ListServiceAccountResponse](#n0stack.iam.v1alpha.ListServiceAccountResponse) |  |
| GetServiceAccount | [GetServiceAccountRequest](#n0stack.iam.v1alpha.GetServiceAccountRequest) | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) |  |
| UpdateServiceAccount | [UpdateServiceAccountRequest](#n0stack.iam.v1alpha.UpdateServiceAccountRequest) | [ServiceAccount](#n0stack.iam.v1alpha.ServiceAccount) |  |
| DeleteServiceAccount | [DeleteServiceAccountRequest](#n0stack.iam.v1alpha.DeleteServiceAccountRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



<a name="n0stack/iam/v1alpha/user.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## n0stack/iam/v1alpha/user.proto



<a name="n0stack.iam.v1alpha.CreateUserRequest"></a>

### CreateUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#n0stack.iam.v1alpha.User) |  |  |






<a name="n0stack.iam.v1alpha.DeleteUserRequest"></a>

### DeleteUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.UpdateUserRequest"></a>

### UpdateUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#n0stack.iam.v1alpha.User) |  |  |
| update_mask | [google.protobuf.FieldMask](#google.protobuf.FieldMask) |  |  |






<a name="n0stack.iam.v1alpha.User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is a unique field. |
| annotations | [User.AnnotationsEntry](#n0stack.iam.v1alpha.User.AnnotationsEntry) | repeated | Annotations can store metadata used by the system for control. In particular, implementation-dependent fields that can not be set as protobuf fields are targeted. The control specified by n0stack may delete metadata specified by the user. |
| labels | [User.LabelsEntry](#n0stack.iam.v1alpha.User.LabelsEntry) | repeated | Labels stores user-defined metadata. The n0stack system must not rewrite this value. |
| display_name | [string](#string) |  |  |
| public_keys | [User.PublicKeysEntry](#n0stack.iam.v1alpha.User.PublicKeysEntry) | repeated |  |






<a name="n0stack.iam.v1alpha.User.AnnotationsEntry"></a>

### User.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.User.LabelsEntry"></a>

### User.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="n0stack.iam.v1alpha.User.PublicKeysEntry"></a>

### User.PublicKeysEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 

 

 


<a name="n0stack.iam.v1alpha.UserService"></a>

### UserService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateUser | [CreateUserRequest](#n0stack.iam.v1alpha.CreateUserRequest) | [User](#n0stack.iam.v1alpha.User) |  |
| GetUser | [GetUserRequest](#n0stack.iam.v1alpha.GetUserRequest) | [User](#n0stack.iam.v1alpha.User) |  |
| UpdateUser | [UpdateUserRequest](#n0stack.iam.v1alpha.UpdateUserRequest) | [User](#n0stack.iam.v1alpha.User) |  |
| DeleteUser | [DeleteUserRequest](#n0stack.iam.v1alpha.DeleteUserRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



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


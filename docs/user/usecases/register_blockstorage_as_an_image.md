# Register blockstorage as an Image

You can manage blockstorages by registering to image, versioning blockstorage with tag.

```yaml
FetchISO:
  type: BlockStorage
  action: FetchBlockStorage
  args:
    name: cloudimage-ubuntu-1804
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824 # 1GiB
    limit_bytes: 10737418240 # 10GiB
    source_url: https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img

ApplyImage:
  type: Image
  action: ApplyImage
  args:
    name: cloudimage-ubuntu

RegisterBlockStorage:
  type: Image
  action: RegisterBlockStorage
  args:
    image_name: cloudimage-ubuntu
    block_storage_name: cloudimage-ubuntu-1804
    tags:
      - latest
      - "18.04"
  depends_on:
    - ApplyImage
```

## Generate BlockStorage from Image

```yaml
GenerateBlockStorage:
  type: Image
  action: GenerateBlockStorage
  args:
    image_name: cloudimage-ubuntu
    tag: "18.04"
    block_storage_name: test-blockstorage
    annotations:
      n0core/provisioning/block_storage/request_node_name: vm-host1
    request_bytes: 1073741824
    limit_bytes: 10737418240
```

## Delete image

```yaml
Remove_cloudimage-ubuntu:
  type: Image
  action: DeleteImage
  args:
    name: cloudimage-ubuntu
  depends_on:
    - Delete_test-vm
```

## Delete image (detailed)

```yaml
Untag_1804_from_cloudimage-ubuntu:
  type: Image
  action: UntagImage
  args:
    name: cloudimage-ubuntu
    tag: "18.04"
  depends_on:
    - Delete_test-vm

Untag_latest_from_cloudimage-ubuntu:
  type: Image
  action: UntagImage
  args:
    name: cloudimage-ubuntu
    tag: latest
  depends_on:
    - Delete_test-vm

Unregister_cloudimage-ubuntu-1804-from-cloudimage-ubuntu:
  type: Image
  action: UnregisterBlockStorage
  args:
    image_name: cloudimage-ubuntu:
    block_storage_name: cloudimage-ubuntu-1804
  depends_on:
    - Untag_1804_from_cloudimage-ubuntu
    - Untag_latest_from_cloudimage-ubuntu

Remove_cloudimage-ubuntu:
  type: Image
  action: DeleteImage
  args:
    name: cloudimage-ubuntu
  depends_on:
    - Unregister_cloudimage-ubuntu-1804-from-cloudimage-ubuntu

Remove_cloudimage-ubuntu-1804:
  type: BlockStorage
  action: DeleteBlockStorage
  args:
    name: cloudimage-ubuntu-1804
  depends_on:
    - Unregister_cloudimage-ubuntu-1804-from-cloudimage-ubuntu
```

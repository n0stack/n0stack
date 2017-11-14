# Graph

オブジェクトは他のオブジェクトとの関係性、グラフで表すことができる。

Objects can explained with relations and graph with other objects.

## Directed graph

オブジェクトを構成するために必要な関係を有向グラフで示す。

Directed graph explain relatioinships necessary to construct objects.

### 懸念点

- 循環参照を起こさないように設計、もしくは安全装置を考える必要がある

### Example 1: Scheduled resource

リソースはスケジューリングされたエージェントの情報が必要である。

Resources need scheduled agent information.

```
(resource/vm/kvm) -[r: n0stack.jp/n0core/scheduled]-> (agent/compute/kvm)
```

### Example 2: Depending resource

他のリソースに依存しているVMやポートのようなオブジェクトは、ボリュームやネットワークなどのリソースの情報が必要である。

Depending resources, like VMs, need some resources informations, like volumes and ports.

```
(resource/vm/kvm) -[r: n0stack.jp/n0core/vm/attachment, n0stack.jp/resource/vm/boot_priority: 1]-> (resource/volume/file)
                  -[r: n0stack.jp/n0core/vm/attachment]-> (resource/port) -[n0stack.jp/resource/port/network: true]-> (resource/network/vlan)
```

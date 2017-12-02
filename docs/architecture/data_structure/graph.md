# Graph

オブジェクトは他のオブジェクトとの関係性をグラフで表すことができる。

## Directed graph

オブジェクトを構成するために必要な関係を有向グラフで示す。依存している理由を `r` キーで示す。

### 懸念点

- 循環参照を起こさないように設計、もしくは安全装置を考える必要がある

### Example 1: Scheduled resource

リソースはスケジューリングされたエージェントの情報が必要である。

```
(resource/vm/kvm) -[r: n0stack/n0core/scheduled]-> (agent/compute/kvm)
```

### Example 2: Depending resource

他のリソースに依存しているVMやポートのようなオブジェクトは、ボリュームやネットワークなどのリソースの情報が必要である。

```
(resource/vm/kvm) -[r: n0stack/n0core/vm/attachment, n0stack/resource/vm/boot_priority: 1]-> (resource/volume/file)
                  -[r: n0stack/n0core/vm/attachment]-> (resource/nic) -[n0stack/resource/nic/network: true]-> (resource/network/vlan)
```

# Compute

## Principle

- Agentにはシンプルな命令方式のgRPCインターフェイス( kvm, iproute2 )を実装
- `ApplyCompute` で適用
- ステートについてはgoroutineでqmpのイベントを監視し、Componentで更新
  - `annotations.n0core/compute/state` にステートをいれる

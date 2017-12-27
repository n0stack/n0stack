# Architecture

n0coreはresource(VM, volume, network)を管理するためのコンポーネントである。

## 0. Overview

絶対書く

## 1. Data structures

### 1.1 [Models](data_structure/models.md)

モデルは互いの関係を[有向プロパティグラフ](data_structure/graph.md)構造で表すことができる。

- [Resources](data_structure/resources.md)
  - VM
  - Volume
  - Network
  - NIC
- [Agents](data_structure/agents.md)
  - Compute Agent
  - Volume Agent
  - Network Agent

### 1.2 [Messages](data_structure/__init__.md)

MessageはModelを１つまたは複数他のProcessorに通知を行う

![](messages.png)

- [Spec](/n0core/message/spec.py)
- [Notify](/n0core/message/notify.py)

## 2. Data flow

過去の名前が残っており、以下の対応である。ある程度出来たら修正する予定。

```
Scheduler == Distributor
Conductor == Aggregater
```

- [Clinet -> API](data_flow/client2api.md)
- [API -> Distributor](data_flow/api2distributor.md)
- [Distributor -> agents -> Aggregater](data_flow/agent.md)

## 3. Architecture overview (similar to onion)

Sphinxのドキュメントから生成する

### 3.1 [Processor](/n0core/processor/__init__.py) (Domain service)

抽象的な処理を書く

- [API]()
- [Distributor](/n0core/processor/distributor.py)
- [Aggregater](/n0core/processor/aggregater.py)
- [Agent](/n0core/processor/agent.py)

### 3.2 Application service

各種フレームワークを使う目的は以下のように３種類にわけられると考えている。

- [Gateway](/n0core/gateway/__init__.py): データの流れ
  - エンドユーザーや他のProcessorへの通信をMQやHTTPなどを使っておこなう
- [Target](/n0core/target/__init__.py): リソースの展開
  - Modelを渡すとKVMやiproute2のようなサービスを使ってリソースを展開する
- [Repository](/n0core/repository/__init__.py): データの永続化
  - RDBMSやGraphDBを使ってProcessorが使うようなインターフェイスを提供する

## Appendix

- [Example: launching VM](ex_launching_vm.md)

# Architecture

n0coreはresource(VM, volume, network)を管理するためのコンポーネントである。

N0core is a component to manage computer resources like memory, cpu, volume and network.

## 0. Overview

絶対書く

## 1. Data structures

- [Graph](data_structure/graph.md)
- [Objects](data_structure/objects.md)

### 1.1 Fact data

- [Resources](data_structure/resources.md)
- [Agents](data_structure/agents.md)

### 1.2 Storing data

- [RDBMS](data_structure/rdbms.md)
- [GraphDB](data_structure/graphdb.md)

## 2. Data flow

- [Messages](data_flow/messages.md)
- [Clinet -> API](data_flow/client2api.md)
- [API -> Schedulers](data_flow/api2scheduler.md)
- [Schedulers -> agents -> Conductor](data_flow/agent.md)

## 3. Procceses overview

これから

### 3.1 Workers

- [API](procceses/api.md)
- [scheduler](procceses/scheduler.md)
- [conductor](procceses/conductor.md)

### 3.2 Agents

- [compute](procceses/compute.md)
- [volumer](procceses/volumer.md)
- [porter](procceses/porter.md)

## Appendix

- [Example: launching VM](ex_launching_vm.md)

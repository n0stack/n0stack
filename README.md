# n0core

[![Build Status](https://travis-ci.org/n0stack/n0core.svg?branch=develop)](https://travis-ci.org/n0stack/n0core)

## Motivation

- 物理的なリソースを仮想的に使うようにするためのものである
- 本リポジトリは他のコンポーネントを開発するためのフレームワークを作成していくための雛形である

## Environment

- Ubuntu 16.04 LTS (Xenial Xerus)
- Golang 1.9

## Dependencies

### kvm

- qemu-kvm

### tap

- iproute2

### qcow2

- qemu-utils

## How to run

### Install packages

```sh
sudo apt install -y \
  iproute2 \
  qemu-kvm \
  qemu-utils
```

### Install libraries

- `proto.go` の更新が早いため定期的にやってほしい

```sh
dep ensure -update
```

### Run n0core

- `sudo go run main.go`
- build and run binary with `sudo`.

## 構成

- Agent
- API

#### 各実装は各ディレクトリの `README.md` を参照のこと

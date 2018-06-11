# n0core

[![Build Status](https://travis-ci.org/n0stack/n0core.svg?branch=develop)](https://travis-ci.org/n0stack/n0core)

## Motivation

- 物理的なリソースを仮想的に使うようにするためのものである
- 本リポジトリは他のコンポーネントを開発するためのフレームワークを作成していくための雛形である

## Principle

- n0coreが死んでも、壊れてもサービス(データプレーン)に影響がないようにする
- 構成ファイルのバックアップがある限り、すでにあるデータプレーンに適合してリストアできるようにする

## Environment

- Ubuntu 16.04 LTS (Xenial Xerus)
- Golang 1.10

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
dep ensure
dep prune
```

### Run on local

```sh
docker-compose up --build etcd api
sudo go run cmd/agent/main.go serve --name=test-node --advertise-address=.. --api-address=`docker inspect -f '{{.NetworkSettings.Networks.n0core_default.IPAddress}}' n0core_api_1` --api-port=20180
```

## 構成

- Agent
- API

#### 各実装は各ディレクトリの `README.md` を参照のこと

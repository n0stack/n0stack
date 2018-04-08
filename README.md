# n0core

[![Build Status](https://travis-ci.org/n0stack/n0core.svg?branch=develop)](https://travis-ci.org/n0stack/n0core)

## Motivation

### n0stack

- クラウドクラスタ全体をマネージできるものを構築する
- Simple and small system
  - 設定を四苦八苦するよりもわかりやすくて短いソースを読んだほうが理解できるし、何なら自分で書いたほうがいい
- 新しい技術や論文を積極的に取り入れて楽しい物を作る

### n0core

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
dep ensure
```

### Run n0core

- `sudo go run main.go`
- build and run binary with `sudo`.

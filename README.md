# n0core

[![Build Status](https://travis-ci.org/n0stack/n0core.svg?branch=master)](https://travis-ci.org/n0stack/n0core)

## Motivation

- 物理的なリソースを仮想的に使うようにするためのものである
- 本リポジトリは他のコンポーネントを開発するためのフレームワークを作成していくための雛形である

## Principle

- n0coreが死んでも、壊れてもサービス(データプレーン)に影響がないようにする
- 構成ファイルのバックアップがある限り、すでにあるデータプレーンに適合してリストアできるようにする

## Environment

- Ubuntu 16.04 LTS (Xenial Xerus)
- Golang 1.10

## How to develop

### Install libraries

```sh
make dep
```

### Tests

#### small

- only localhost
- short time

```sh
make test-small
make test-small-v
make test-small-docker
```

#### medium

- with root
- having dependency for outside

```sh
make test-medium
make test-medium-v
```

#### 各実装は各ディレクトリの `README.md` を参照のこと

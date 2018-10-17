# n0core

[![Build Status](https://travis-ci.org/n0stack/n0core.svg?branch=master)](https://travis-ci.org/n0stack/n0core)

## Motivation

- The example for implementation of n0stack API
- 本リポジトリは他のコンポーネントを開発するための雛形である

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

#### Upgrade libraries

```sh
make dep-update
```

### Build

```sh
make build
make build-docker
make build-proto
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
make test-medium-without-root
make test-medium-without-external
```

### Run all in one

```sh
make run-all-in-one
```

## Dependency map

![](docs/images/dependency_map.svg)

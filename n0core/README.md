# n0core

## Motivation

- The example for implementation of n0stack API
- 本リポジトリは他のコンポーネントを開発するための雛形である

## Principle

- n0coreが死んでも、壊れてもサービス(データプレーン)に影響がないようにする
- 構成ファイルのバックアップがある限り、すでにあるデータプレーンに適合してリストアできるようにする

## Environment

- Ubuntu 18.04 LTS (Bionic Beaver)
- Golang 1.10

## How to deploy

### API

```
cd ..
make up
```

### Agent

#### Remote

- rootでログインできる必要がある
- sftp でバイナリを送り、 systemd サービスを起動する
- ファイルは `/var/lib/n0core` に送られ、シンボリックリンクが貼られる

```
bin/n0core deploy agent -i id_ecdsa root@$node_id -name vm-host1 -advertise-address=$node_id -node-api-endpoint=$api_address:20180
```

#### Local

```
bin/n0core install agent -a "-name vm-host1 -advertise-address=$node_id -node-api-endpoint=$api_address:20180"
```

## Dependency map

![](../docs/images/dependency_map.svg)

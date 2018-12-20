# n0cli

CLI for n0stack API.

<!-- ## Demo -->

## Usage

See also command help.

```
% bin/n0cli -h
NAME:
   n0cli - the n0stack CLI application

USAGE:
   n0cli [global options] command [command options] [arguments...]

VERSION:
   28

COMMANDS:
     get      Get resource if set resource name, List resources if not set
     delete   Delete resource
     do       Do DAG tasks (Detail n0stack/pkg/dag)
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-endpoint value  (default: "localhost:20180") [$N0CLI_API_ENDPOINT]
   --help, -h            show help
   --version, -v         print the version

```

```
% bin/n0cli get -h
NAME:
   n0cli get - Get resource if set resource name, List resources if not set

USAGE:
   n0cli get [resource type] [resource name]
```

```
% bin/n0cli delete -h
NAME:
   n0cli delete - Delete resource

USAGE:
   n0cli delete [resource type] [resource name]

% bin/n0cli do -h
NAME:
   n0cli do - Do DAG tasks (Detail n0stack/pkg/dag)

USAGE:
   n0cli do [file name]

DESCRIPTION:

  ## File format

  ---
  task_name:
    type: Network
    action: GetNetwork
    args:
      name: test-network
    depend_on:
      - dependency_task_name
    ignore_error: true
  dependency_task_name:
    type: ...
  ---

  - task_name
      - 任意の名前をつけ、ひとつのリクエストに対してユニークなものにする
  - type
      - gRPC メッセージを指定する
      - VirtualMachine や virtual_machine という形で指定できる
  - action
      - gRPC の RPC を指定する
      - GetNetwork など定義のとおりに書く
  - args
      - gRPC の RPCのリクエストを書く
  - depend_on
      - DAG スケジューリングに用いられる
      - task_name を指定する
  - ignore_error
      - タスクでエラーが発生しても継続する
```

## Environment

- Ubuntu 18.04 LTS (Bionic Beaver)

## How to build

```
cd ..
make build-n0cli
```

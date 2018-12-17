# n0core

The example for implementation of n0stack API.

## Environment

- Ubuntu 18.04 LTS (Bionic Beaver)
- Golang 1.11

## How to deploy

### API

```
cd ..
make up
```

### Agent

#### Remote

- Require root user
- Perform the following processing
    - Send self to `/var/lib/n0core/n0core.$VERSION` with sftp
    - Run `n0core local`

```
bin/n0core deploy agent -i id_ecdsa root@$node_ip -name vm-host1 -advertise-address=$node_ip -node-api-endpoint=$api_address:20180
```

#### Local

- Require root user
- Perform the following processing
    - If n0core service is started, stop n0core service.
    - Create symbolic link from self to `/usr/bin/n0core`
    - Generate systemd unit file and start systemd service

```
bin/n0core install agent -a "-name vm-host1 -advertise-address=$node_ip -node-api-endpoint=$api_address:20180"
```

## How to develop

- see also [Makefile](../Makefile)

### Build

```
cd ..
make build-n0core
```

# n0stack

[![Build Status](https://travis-ci.org/n0stack/n0stack.svg?branch=master)](https://travis-ci.org/n0stack/n0stack)
[![CircleCI](https://circleci.com/gh/n0stack/n0stack/tree/master.svg?style=shield)](https://circleci.com/gh/n0stack/n0stack/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/n0stack/n0stack)](https://goreportcard.com/report/github.com/n0stack/n0stack)
[![](https://img.shields.io/docker/pulls/n0stack/n0stack.svg)](https://hub.docker.com/r/n0stack/n0stack)
<!-- [![](https://img.shields.io/docker/build/n0stack/n0stack.svg)](https://hub.docker.com/r/n0stack/n0stack) -->

The n0stack is a simple cloud provider using gRPC.

## Description

The n0stack is...

- a cloud provider.
    - You can use some features: booting VMs, managing networks and so on (see also [n0proto](n0proto/).)
- simple.
    - There are shortcode and fewer options.
- using gRPC.
    - A unified interface increase reusability.
- able to be used as library and framework.
    - You can concentrate to develop your logic by sharing libraries and frameworks for middleware, test, and deployment.

## Motivation

Cloud providers have various forms depending on users.
This problem has been solved with many options and add-ons (e.g. OpenStack configuration file is very long.)
It is difficult to adapt to the application with options, therefore it is necessary to read or rewrite long abstracted code.
I think it is better to code it yourself from the beginning.

There are some problems to develop cloud providers from scratch: no libraries, software quality, man-hours, and deployment.
The n0stack wants to solve these problems.

<!-- ## Demo -->

## Getting started

### Prerequisites

- Docker
- docker-compose
- Ubuntu 18.04 LTS

### Deploy all in one

1. You can start controllers on docker and install agent as follows:

```sh
wget https://raw.githubusercontent.com/n0stack/n0stack/master/deploy/docker-compose.yml
docker-compose up -d
docker run -it --rm -v $PWD:/dst n0stack/n0stack cp /usr/local/bin/n0core /dst
./n0core install agent -a "--node-api-endpoint=localhost:20180 --location=////1"
```

2. Download n0cli from [Github releases](https://github.com/n0stack/n0stack/releases/latest).
3. Try [use cases](https://docs.n0st.ac/en/master/user/usecases/README.html).

## Documentations

[![Gitter](https://img.shields.io/gitter/room/n0stack/n0sack.svg)](https://gitter.im/n0stack/)
[![Documentation Status](https://readthedocs.org/projects/n0stack/badge/?version=master)](https://docs.n0st.ac/en/master/?badge=master)
[![GoDoc](https://godoc.org/github.com/n0stack/n0stack?status.svg)](https://godoc.org/github.com/n0stack/n0stack)

User documentations and specifications is [readthedocs](https://docs.n0st.ac/en/master/?badge=master).

Golang library documentations is [GoDoc](https://godoc.org/github.com/n0stack/n0stack).

## Components

The final goal of n0stack is to represent the state of all clusters with n0proto.
Implementations such as n0core manipulates the cluster according to the information specified by n0proto.
The implementation of n0proto is left to each developer.
This repository is just a reference implementation.
However, please share actively usable libraries such as `n0core/pkg/driver`.

![](docs/_static/images/components.svg)

### [n0proto](n0proto/)

Protobuf definitions for all of n0stack services.

### [n0cli](n0cli/)

CLI for n0stack API.

### n0ui

Web UI for n0stack API.

### [n0bff](n0bff/)

BFF(Backends for Frontend) of n0stack API. This provide features: API gateway, authentication, authorization and so on.

### [n0core](n0core/)

The example for implementations about n0stack API.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request 

## License

[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fn0stack%2Fn0stack.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fn0stack%2Fn0stack?ref=badge_large)

# n0stack

Builds: 
[![Build Status](https://travis-ci.org/n0stack/n0stack.svg?branch=master)](https://travis-ci.org/n0stack/n0stack)
[![CircleCI](https://circleci.com/gh/n0stack/n0stack/tree/master.svg?style=shield)](https://circleci.com/gh/n0stack/n0stack/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/n0stack/n0stack)](https://goreportcard.com/report/github.com/n0stack/n0stack)
[![](https://img.shields.io/docker/pulls/n0stack/n0stack.svg)](https://hub.docker.com/r/n0stack/n0stack)
<!-- [![](https://img.shields.io/docker/build/n0stack/n0stack.svg)](https://hub.docker.com/r/n0stack/n0stack) -->

Documentations: 
[![Gitter](https://img.shields.io/gitter/room/n0stack/n0sack.svg)](https://gitter.im/n0stack/)
[![Documentation Status](https://readthedocs.org/projects/n0stack/badge/?version=master)](https://docs.n0st.ac/en/master/?badge=master)
[![GoDoc](https://godoc.org/github.com/n0stack/n0stack?status.svg)](https://godoc.org/github.com/n0stack/n0stack)

License: 
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fn0stack%2Fn0stack.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fn0stack%2Fn0stack?ref=badge_shield)

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
However, it is difficult to adapt to the application by options, then it is necessary to read or rewrite long abstracted codes.
Therefore, I thought that it would be better to code on your hands from beginning.

There are some problems to develop cloud providers from scratch: no library, software quality, man-hour, and deployment.
The n0stack wants to solve such problems.

<!-- ## Demo -->

## Components

![](docs/_static/images/components.svg)

### [n0proto](n0proto/)

Protobuf definitions for all of n0stack services.

### [n0cli](n0cli/)

CLI for n0stack API.

<!-- ### n0ui

Web UI for n0stack API -->

### [n0core](n0core/)

The example for implementations about n0stack API.

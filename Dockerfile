FROM golang:1.13 AS BUILD_GO

ENV GO111MODULE=on

COPY . /go/src/github.com/n0stack/n0stack
WORKDIR /go/src/github.com/n0stack/n0stack

RUN make build-go

FROM debian:jessie

RUN apt update \
 && apt install -y openssh-client \
 && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

COPY VERSION /
COPY LICENSE /
COPY --from=BUILD_GO /go/src/github.com/n0stack/n0stack/bin/* /usr/local/bin/

WORKDIR /root
CMD /bin/bash

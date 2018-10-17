FROM golang:1.11 AS build

RUN apt-get update \
 && apt-get install -y \
        cloud-image-utils \
        iproute2 \
        qemu-kvm \
        qemu-utils \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/n0stack/n0stack/n0core
COPY Gopkg.toml Gopkg.lock ./

RUN go get -u github.com/golang/dep/cmd/dep \
 && dep ensure -v -vendor-only=true

RUN go get -u golang.org/x/lint/golint

COPY . /go/src/github.com/n0stack/n0stack/n0core

ENV DISABLE_KVM=1

# 遅い
# RUN make test-small
RUN make build

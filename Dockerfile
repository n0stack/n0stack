FROM golang:1.11 AS build

RUN apt-get update \
 && apt-get install -y \
        qemu-kvm \
        qemu-utils \
        iproute2 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/n0stack/n0core
COPY Gopkg.toml Gopkg.lock ./

RUN go get -u github.com/golang/dep/cmd/dep \
 && dep ensure -v -vendor-only=true

RUN go get -u golang.org/x/lint/golint

COPY . /go/src/github.com/n0stack/n0core

# RUN make analysis \
#  && make test-small


# RUN go build -o /api cmd/api/main.go \
#  && go build -o /agent cmd/agent/main.go

# FROM debian:jessie

# COPY --from=build /api /api
# COPY --from=build /agent /agent

# WORKDIR /

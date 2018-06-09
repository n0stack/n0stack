FROM golang:1.10 AS build

WORKDIR /go/src/github.com/n0stack/n0core
COPY . /go/src/github.com/n0stack/n0core

RUN go build -o /api cmd/api/main.go \
 && go build -o /agent cmd/agent/main.go

FROM debian:jessie

COPY --from=build /api /api
COPY --from=build /agent /agent

WORKDIR /

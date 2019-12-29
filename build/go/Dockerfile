FROM golang:1.13

# for n0core
RUN apt-get update \
 && apt-get install -y \
        cloud-image-utils \
        iproute2 \
        qemu-kvm \
        qemu-utils \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN go get -u golang.org/x/lint/golint

ENV GO111MODULE=on
ENV DISABLE_KVM=1

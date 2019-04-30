FROM golang:1.11 AS BUILD
LABEL maintainer="h-otter@outlook.jp"

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
 && go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
 && go get -u github.com/golang/protobuf/protoc-gen-go

RUN apt update \
 && apt install -y unzip \
 && cd /tmp \
 && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip \
 && unzip protoc-3.7.0-linux-x86_64.zip \
 && mv bin/protoc /usr/bin/

 RUN cd /tmp \
  && wget https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 \
  && chmod 755 swagger_linux_amd64 \
  && mv /tmp/swagger_linux_amd64 /usr/bin/swagger

WORKDIR /src
COPY entry_point.sh /
COPY swagger.sh /

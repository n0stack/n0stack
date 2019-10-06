FROM pseudomuto/protoc-gen-doc AS GEN_DOC

FROM golang:1.13 AS BUILD
LABEL maintainer="h-otter@outlook.jp"

COPY --from=GEN_DOC /usr/local/bin/protoc-gen-doc /usr/bin/protoc-gen-doc

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
 && go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
 && go get -u github.com/golang/protobuf/protoc-gen-go

RUN apt update \
 && apt install -y unzip \
 && cd /tmp \
 && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip \
 && unzip protoc-3.7.0-linux-x86_64.zip \
 && mv bin/protoc /usr/bin/ \
 && chmod 755 /usr/bin/protoc

RUN cd /tmp \
  && wget https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 \
  && chmod 755 swagger_linux_amd64 \
  && mv /tmp/swagger_linux_amd64 /usr/bin/swagger

RUN cp -r /go/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger /tmp/include \
 && cp -r /go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google /tmp/include \
 && chmod -R 755 /tmp/include

WORKDIR /src
COPY entry_point.sh /

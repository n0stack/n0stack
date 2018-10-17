FROM grpc/go
LABEL maintainer="h-otter@outlook.jp"

RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
 && go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

WORKDIR /src
COPY entry_point.sh /

# ENTRYPOINT [ "/entry_point.sh" ]

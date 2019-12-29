FROM n0stack/build-grpc-go AS LIB

FROM python:3.7
LABEL maintainer="h-otter@outlook.jp"

RUN pip install \
      googleapis-common-protos \
      grpcio-tools

COPY --from=LIB /tmp/include /tmp/include

WORKDIR /src
COPY entry_point.sh /

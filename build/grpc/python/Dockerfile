FROM python:3.7
LABEL maintainer="h-otter@outlook.jp"

RUN pip install \
    googleapis-common-protos \
    grpcio-tools

WORKDIR /src
COPY entry_point.sh /

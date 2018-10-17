#!/bin/bash
set -x

dirs=`find /src -type d | grep -v .git | grep -v test`

for d in $dirs
do
  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    protoc \
      -I/usr/local/include \
      -I/src \
      -I${GOPATH}/src \
      -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      $* $d/*.proto
  fi
done

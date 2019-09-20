#!/bin/bash
set -x

dirs=`find /src -type d | egrep -v "/\." | grep -v test | grep -v vendor`

for d in $dirs
do
  rm $d/*.pb.go
  rm $d/*.pb.gw.go

  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    protoc \
      -I/src \
      -I/tmp/include \
      -I${GOPATH}/src \
      -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
      -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      $* $d/*.proto
  fi
done

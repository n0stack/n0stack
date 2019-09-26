#!/bin/bash
set -x

dirs=`find /src -type d | egrep -v "/\." | grep -v test | grep -v vendor`
echo "{}" > $*/n0stack.swagger.json

for d in $dirs
do
  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    name=${d#\/src\/}

    protoc \
      -I/src \
      -I/tmp/include \
      -I${GOPATH}/src \
      -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
      -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --doc_out=/dst \
      --doc_opt=markdown,${name//\//_}.md
      $d/*.proto

    # if [ "$?" != "0" ]; then
    #   exit 1
    # fi
  fi
done

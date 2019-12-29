#!/bin/bash
set -x

dirs=`find /src -type d | egrep -v "/\." | grep -v test | grep -v vendor`

for d in $dirs
do
  # touch $d/__init__.py

  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    python \
      -m grpc_tools.protoc \
      -I/src \
      -I/tmp/include \
      --python_out=/src \
      --grpc_python_out=/src \
      $* $d/*.proto
  fi
done

#!/bin/bash
set -x

dirs=`find /src -type d | grep -v .git | grep -v test`

for d in $dirs
do
  # rm $d/*.py
  # touch $d/__init__.py

  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    python \
      -m grpc_tools.protoc \
      -I/usr/local/include \
      -I/src \
      --python_out=/dst \
      --grpc_python_out=/dst \
      $* $d/*.proto
  fi
done

#!/bin/bash
set -x

mkdir -p /tmp/build/n0proto
mkdir -p /tmp/dst
cp -r /src/* /tmp/build/n0proto
dirs=`find /tmp/build/n0proto -type d | grep -v .git | grep -v test`
rm -r /dst/*

for d in $dirs
do
  # touch $d/__init__.py

  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    python \
      -m grpc_tools.protoc \
      -I/usr/local/include \
      -I/tmp/build/n0proto \
      --python_out=/tmp/dst \
      --grpc_python_out=/tmp/dst \
      $* $d/*.proto
  fi
done

mv /tmp/dst/n0proto/* /dst

dirs=`find /dst -type d | grep -v .git | grep -v test`
for d in $dirs
do
  touch $d/__init__.py
done

#!/bin/bash
set -x

dirs=`find /src -type d | egrep -v "/\." | grep -v test | grep -v vendor`

for d in $dirs
do
  rm -f $d/*.pb.go
  rm -f $d/*.pb.gw.go

  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    protoc \
      -I/src \
      -I/tmp/include \
      $* $d/*.proto

    # if [ "$?" != "0" ]; then
    #   exit 1
    # fi
  fi
done

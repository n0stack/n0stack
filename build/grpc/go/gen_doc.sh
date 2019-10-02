#!/bin/bash
set -x

dirs=`find /src -type d | egrep -v "/\." | grep -v test | grep -v vendor`

for d in $dirs
do
  # 複数のファイルを指定できない
  ls -1 $d/*.proto > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    name=${d#\/src\/}
    protoc \
      -I/src \
      -I/tmp/include \
      --doc_out=/dst \
      --doc_opt=markdown,${name//\//_}.md \
      $d/*.proto

    # if [ "$?" != "0" ]; then
    #   exit 1
    # fi
  fi
done

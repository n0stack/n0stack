#!/bin/bash
set -x
set -e

rm -f /doc_dst/*

for wd in `find /src -maxdepth 1 -mindepth 1 -type d`
do
  dirs=`find $wd -type d | egrep -v "/\." | grep -v test | grep -v vendor | grep -v sandbox`
  echo "{}" > $wd/`basename $wd`.swagger.json

  for d in $dirs
  do
    rm -f $d/*.pb.go
    rm -f $d/*.pb.gw.go

    if ls $d/*.proto > /dev/null 2>&1
    then
      name=${d#\/src\/}
      protoc \
        -I/src \
        -I/tmp/include \
        --go_out=plugins=grpc:/go/src \
        --grpc-gateway_out=logtostderr=true:/go/src \
        --doc_out=/doc_dst \
        --doc_opt=markdown,${name//\//_}.md \
        --swagger_out=logtostderr=true,allow_merge=true,merge_file_name=n0stack:$wd \
        $d/*.proto
    fi
  done
done

# build-grpc-py

## how to build

```sh
docker build -t n0stack/build-gprc-py build/grpc/python
docker run -it --rm -v $PWD/n0proto:/src n0stack/build-gprc-py /entry_point.sh
```

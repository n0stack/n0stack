FROM n0stack/build-go AS BUILD_GO

COPY . /go/src/github.com/n0stack/n0stack
WORKDIR /go/src/github.com/n0stack/n0stack

RUN make test-small-go \
 && make build-go

FROM debian:jessie

COPY VERSION /
COPY LICENSE /
COPY --from=BUILD_GO bin/* /usr/bin/

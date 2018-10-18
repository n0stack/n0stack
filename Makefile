GOOS=linux
GOARCH=amd64
GOCMD=go


.PHONY: build-docker
build-docker:
	docker build -t n0stack/n0core n0core
	docker build -t n0stack/build-grpc-go build/grpc/go
	docker build -t n0stack/build-grpc-py build/grpc/python
.PHONY: build-proto
build-proto:
	docker run -it --rm -v $(PWD)/n0proto:/src:ro -v `go env GOPATH`/src:/dst n0stack/build-proto /entry_point.sh --go_out=plugins=grpc:/dst
	docker run -it --rm -v $(PWD)/n0proto:/src n0stack/build-gprc-py /entry_point.sh
.PHONY: build
build:
	docker run -it --rm -v $(PWD)/n0core:/go/src/github.com/n0stack/n0stack/n0core n0stack/n0core make build  # よしなに修正予定

# test-small:
# 	cd n0core
# 	make test-small
# test-medium:
# 	cd n0core
# 	make test-medium

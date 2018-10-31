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
	go build -o bin/n0core -v ./n0core/cmd/n0core
	go build -o bin/n0cli -v ./n0core/cmd/n0cli
.PHONY: build-on-docker
build-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \ # goのcacheを有効化するために必要
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		n0stack/build-go make build

clean:
	go clean
	sudo rm -rf .go-build
	sudo rm -rf bin
	sudo rm -rf sandbox

test-small:
	go test -cover ./...

# test-medium: up-mock # with root, having dependency for external
# 	sudo go test -tags=medium -cover ./...
# test-medium-v: up-mock
# 	sudo go test -tags=medium -v -cover ./...
# test-medium-without-root: up-mock
# 	go test -tags="medium without_root" -cover ./...
# test-medium-without-external:
# 	sudo go test -tags="medium without_external" -cover ./...

# test-small:
# 	cd n0core
# 	make test-small
# test-medium:
# 	cd n0core
# 	make test-medium

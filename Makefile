GOOS=linux
GOARCH=amd64
GOCMD=go


# --- Deployment ---
run-all-in-one: build-on-docker up
	sudo bin/n0core agent \
		--name=run-all-in-one \
		--advertise-address=10.20.180.1 \
		--node-api-endpoint=localhost:20180 \
		--base-directory=./sandbox/workdir

.PHONY: up
up:
	mkdir -p sandbox
	docker-compose up -d --scale mock_agent=0


# --- Build ---
.PHONY: all
all: build-builder vendor-on-docker build-n0proto-on-docker build-n0core-on-docker build-n0cli-on-docker

.PHONY: build-n0core
build-n0core:
	go build -o bin/n0core -v ./n0core/cmd/n0core

.PHONY: build-n0core-on-docker
build-n0core-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=off \
		n0stack/build-go \
			make build-n0core

.PHONY: build-n0cli
build-n0cli:
	go build -o bin/n0cli -v ./n0core/cmd/n0cli

.PHONY: build-n0cli-on-docker
build-n0cli-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=off \
		n0stack/build-go \
			make build-n0cli

.PHONY: build-builder
build-builder:
	docker build -t n0stack/build-grpc-go build/grpc/go
	docker build -t n0stack/build-grpc-py build/grpc/python
	docker build -t n0stack/build-go build/go

.PHONY: build-n0proto-on-docker
build-n0proto-on-docker:
	docker run -it --rm \
		-v $(PWD)/n0proto:/src:ro \
		-v `go env GOPATH`/src:/dst \
		n0stack/build-grpc-go \
			/entry_point.sh --go_out=plugins=grpc:/dst
	docker run -it --rm \
		-v $(PWD)/n0proto:/src \
		n0stack/build-grpc-py \
			/entry_point.sh


# -- Maintenance ---
.PHONY: update
update: update-go update-novnc

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: vendor-on-docker
vendor-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=on \
		n0stack/build-go \
			make vendor

.PHONY: update-go
update-go:
	go get -u

.PHONY: update-novnc
update-novnc:
	go get -v github.com/rakyll/statik
	rm -rf /tmp/novnc
	git clone --depth 1 https://github.com/novnc/noVNC /tmp/novnc
	statik -p provisioning -Z -f -src /tmp/novnc -dest pkg/api

.PHONY: clean
clean:
	# go clean
	sudo rm -rf .go-build
	sudo rm -rf bin
	sudo rm -rf sandbox
	sudo rm -rf vendor

up-mock:
	mkdir -p sandbox
	docker-compose up -d


# --- Test ---
analysis:
	gofmt -d -s `find ./ -name "*.go" | grep -v vendor`
	golint ./... | grep -v vendor # https://github.com/golang/lint/issues/320

# TODO: check n0proto changes
test-small: build-n0proto-on-docker
	git diff --name-status --exit-code n0proto  # n0proto
	go test -cover ./...  # n0core, n0cli

test-small-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=off \
		n0stack/build-go \
			make test-small

test-medium: build-n0core-on-docker up-mock # with root, having dependency for external
	sudo go test -tags=medium -cover ./...   # n0core, n0cli

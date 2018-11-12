GOOS=linux
GOARCH=amd64
GOCMD=go
VERSION=$(shell cat VERSION)


# --- Deployment ---
run-all-in-one: build-n0core-on-docker up
	sudo bin/n0core agent \
		--name=run-all-in-one \
		--advertise-address=10.20.180.1 \
		--node-api-endpoint=localhost:20180 \
		--base-directory=./sandbox/workdir

.PHONY: up
up: build-n0core-on-docker
	mkdir -p sandbox
	docker-compose up -d --scale mock_agent=0
	docker-compose restart api # reload binary

up-mock: build-n0core-on-docker
	mkdir -p sandbox
	docker-compose up -d
	docker-compose restart api # reload binary
	docker-compose restart mock_agent # reload binary


# --- Build ---
.PHONY: all
all: build-builder vendor-on-docker build-n0proto-on-docker build-n0core-on-docker build-n0cli-on-docker

.PHONY: build-n0core
build-n0core:
	go build -o bin/n0core -ldflags "-X main.version=$(VERSION)" -v ./n0core/cmd/n0core

.PHONY: build-n0core-on-docker
build-n0core-on-docker:
	docker run -it --rm \
		-v $(PWD)/n0core:/src:ro \
		-v `go env GOPATH`/src:/dst \
		n0stack/build-grpc-go \
			/entry_point.sh --go_out=plugins=grpc:/dst
	sudo chown -R $(USER) n0core
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=off \
		n0stack/build-go \
			make build-n0core

.PHONY: build-n0cli
build-n0cli:
	go build -o bin/n0cli -ldflags "-X main.version=$(VERSION)" -v ./n0cli/cmd/n0cli

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
	sudo chown -R $(USER) n0proto.go
	docker run -it --rm \
		-v $(PWD)/n0proto:/src:ro \
		-v $(PWD)/n0proto.py:/dst \
		n0stack/build-grpc-py \
			/entry_point.sh
	sudo chown -R $(USER) n0proto.py

.PHONY: build-versioning
build-versioning:
	go build -o bin/versioning -ldflags "-X main.version=$(VERSION)" -v ./build/versioning


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
	sudo chown -R $(USER) vendor

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
	docker-compose down
	# sudo rm -rf .go-build
	sudo rm -rf bin
	sudo rm -rf sandbox
	# sudo rm -rf vendor

logs:
	docker-compose logs -f api mock_agent

versioning:
	bin/versioning increment -write


# --- Test ---
analysis:
	gofmt -d -s `find ./ -name "*.go" | grep -v vendor`
	golint ./... | grep -v vendor # https://github.com/golang/lint/issues/320

test-small: test-small-n0proto test-small-go

test-small-on-docker: test-small-n0proto
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=off \
		n0stack/build-go \
			make test-small-go

test-small-n0proto: build-n0proto-on-docker
	git diff --name-status --exit-code n0proto.py
	git diff --name-status --exit-code n0proto.go

test-small-go:
	go test -cover ./...

test-medium: up-mock # with root, having dependency for external
	sudo go test -tags=medium -cover ./...   # n0core, n0cli

GOOS=linux
GOARCH=amd64
GOCMD=go
VERSION=$(shell cat VERSION)


# --- Deployment ---
.PHONY: up
up: build-n0core-on-docker
	mkdir -p sandbox
	docker-compose up -d --scale mock_agent=0
	docker-compose restart api # reload binary


# --- Build ---
.PHONY: all
all: build-builder vendor-on-docker build-n0proto-on-docker build-n0core-on-docker build-n0cli-on-docker

.PHONY: build-go
build-go: build-n0core build-n0cli

.PHONY: build-n0core
build-n0core:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/n0core -ldflags "-X main.version=$(VERSION)" -v ./n0core/cmd/n0core

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
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/n0cli -ldflags "-X main.version=$(VERSION)" -v ./n0cli/cmd/n0cli

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

.PHONY: build-n0version
build-n0version:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/n0version -ldflags "-X main.version=$(VERSION)" -v ./build/n0version


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

.PHONY: update-go-on-docker
update-go-on-docker:
	docker run -it --rm \
		-v $(PWD)/.go-build:/root/.cache/go-build/ \
		-v $(PWD):/go/src/github.com/n0stack/n0stack \
		-w /go/src/github.com/n0stack/n0stack \
		-e GO111MODULE=on \
		n0stack/build-go \
			make update-go

.PHONY: update-novnc
update-novnc:
	go get -v github.com/rakyll/statik
	rm -rf /tmp/novnc /tmp/novnc-src
	git clone --depth 1 https://github.com/novnc/noVNC /tmp/novnc-src
	mkdir -p /tmp/novnc
	# Copy only required files, ref: https://github.com/novnc/noVNC/blob/master/docs/EMBEDDING.md
	cp -r /tmp/novnc-src/app /tmp/novnc-src/core /tmp/novnc-src/vendor /tmp/novnc-src/*.html /tmp/novnc/
	rm ./n0core/pkg/api/provisioning/virtualmachine/statik.go
	statik -p virtualmachine -Z -f -src /tmp/novnc -dest ./n0core/pkg/api/provisioning
	rm -rf /tmp/novnc /tmp/novnc-src

.PHONY: clean
clean:
	# go clean
	docker-compose down
	# sudo rm -rf .go-build
	# sudo rm -rf bin/*
	sudo rm -rf sandbox/*
	# sudo rm -rf vendor

logs:
	docker-compose logs -f api

increment:
	bin/n0version increment -write


GOFLAGS := -ldflags "-X main.version=${VERSION}" -v

.PHONY: release-to-github
release-to-github: build-artifacts-to-release
	GO111MODULE=off go get -u github.com/tcnksm/ghr
	ghr -username n0stack -repository n0stack -commitish $(shell git rev-parse HEAD) -recreate 0.2.$(VERSION) ./artifacts/

.PHONY: build-artifacts-to-release
build-artifacts-to-release: artifacts/n0core_linux_amd64.tar.gz \
                            artifacts/n0cli_linux_amd64.tar.gz \
                            artifacts/n0cli_darwin_amd64.tar.gz \
                            artifacts/n0cli_freebsd_amd64.tar.gz \
                            artifacts/n0cli_windows_amd64.zip

.PHONY: artifacts/%.tar.gz
artifacts/%.tar.gz:
	$(eval BASENAME := $(subst .tar.gz,,$(notdir $@)))
	$(eval BASENAME_WORDS := $(subst _, ,$(BASENAME)))
	$(eval BINNAME := $(word 1,$(BASENAME_WORDS)))
	$(eval GOOS := $(word 2,$(BASENAME_WORDS)))
	$(eval GOARCH := $(word 3,$(BASENAME_WORDS)))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o artifacts/$(BINNAME) $(GOFLAGS) ./$(BINNAME)/cmd/$(BINNAME)
	cd artifacts && tar czvf $(notdir $@) $(BINNAME) --owner=n0stack:0 --group=n0stack:0
	rm artifacts/$(BINNAME)

# windows
.PHONY: vendor artifacts/%.zip
artifacts/%.zip:
	$(eval BASENAME := $(subst .zip,,$(notdir $@)))
	$(eval BASENAME_WORDS := $(subst _, ,$(BASENAME)))
	$(eval BINNAME := $(word 1,$(BASENAME_WORDS)))
	$(eval GOOS := $(word 2,$(BASENAME_WORDS)))
	$(eval GOARCH := $(word 3,$(BASENAME_WORDS)))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o artifacts/$(BINNAME).exe $(GOFLAGS) ./$(BINNAME)/cmd/$(BINNAME)
	cd artifacts && zip $(notdir $@) $(BINNAME).exe
	rm artifacts/$(BINNAME).exe

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
	# git diff --name-status --exit-code n0proto.py
	git diff --name-status --exit-code n0proto.go

test-small-go:
	go test -race -cover ./...

test-medium: up # with root, having dependency for external
	sudo go test -race -tags=medium -cover ./...   # n0core, n0cli

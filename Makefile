# Go パラメータ
GOOS=linux
GOARCH=amd64
GOCMD=go


dep:
	dep ensure
dep-update:
	dep ensure -update
	dep prune
	dep status

build:
	go build -o bin/n0core-agent -v cmd/agent/*.go
	go build -o bin/n0core-api -v cmd/api/*.go
	# go build -o bin/n0stack -v cmd/n0stack/*.go
build-docker:
	docker build -t n0stack/n0core .

up: build-docker
	mkdir -p sandbox
	docker-compose up -d --scale mock_agent=0
up-mock: build-docker
	mkdir -p sandbox
	docker-compose up -d
logs:
	docker-compose logs -f
rm:
	docker-compose down
	docker-compose rm
clean:
	go clean
	rm -rf bin

analysis:
	gofmt -d -s `find ./ -name "*.go" | grep -v vendor`
	golint ./... | grep -v vendor # https://github.com/golang/lint/issues/320

test-small:
	go test -cover ./...
test-small-v:
	go test -v -cover ./...
test-small-docker:
	docker run -it --rm -v $(PWD):/go/src/github.com/n0stack/n0core n0stack/n0core make test-small

test-medium: up-mock # with root, having dependency for external
	go test -tags=medium -cover ./...
test-medium-v: up-mock
	go test -tags=medium -v -cover ./...
test-medium-without-root: up-mock
	go test -tags="medium without_root" -cover ./...
test-medium-without-external:
	go test -tags="medium without_external" -cover ./...

run-all-in-one: up
	docker run --rm -it -v $(PWD)/bin:/go/src/github.com/n0stack/n0core/bin n0stack/n0core make build
	sudo ./bin/n0core-agent serve \
		--name=test \
		--advertise-address=10.20.180.1 \
		--node-api-endpoint=localhost:20181 \
		--base-directory=./sandbox/workdir

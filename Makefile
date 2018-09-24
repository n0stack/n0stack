# Go パラメータ
GOOS=linux
GOARCH=amd64
GOCMD=go


build:
	go build -o bin/n0core-agent -v cmd/agent/*.go
	go build -o bin/n0core-api -v cmd/api/*.go
build-docker:
	docker build -t n0stack/n0core .

up: build-docker
	docker-compose up -d
logs: 
	docker-compose logs -f
run-agent: up
	docker run --rm -it -v $(PWD)/bin:/go/src/github.com/n0stack/n0core/bin n0stack/n0core make build
	./bin/n0core-agent serve \
		--name=test \
		--advertise-address=10.20.180.1 \
		--node-api-endpoint=localhost:20181 \
		--base-directory=./sandbox/workdir

dep:
	dep ensure
dep-update:
	dep ensure -update
	dep prune
	dep status

analysis:
	gofmt -d -s `find ./ -name "*.go" | grep -v vendor`
	golint ./... | grep -v vendor # https://github.com/golang/lint/issues/320

test-small:
	go test -cover ./...
test-small-v:
	go test -v -cover ./...
test-small-docker:
	docker run -it --rm -v $(PWD):/go/src/github.com/n0stack/n0core n0stack/n0core make test-small

test-medium: # with root, having dependency for outside
	go test -tags=medium -cover ./...
test-medium-v:
	go test -tags=medium -v -cover ./...

clean:
	go clean
	rm -rf bin

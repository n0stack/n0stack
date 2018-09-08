# Go パラメータ
GOOS=linux
GOARCH=amd64
GOCMD=go


run_local_agent:
	docker-compose up --build api etcd
	go run cmd/agent/main.go

build:
	go build cmd/agent/*.go -o bin/agent -v -x
build-docker:
	docker build -t n0stack/n0core .

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
test-small-docker: build-docker
	docker run -it --rm -e DISABLE_KVM=1 n0stack/n0core make test-small

test-medium: # with root, having dependency for outside
	go test -tags=medium -cover ./...
test-medium-v:
	go test -tags=medium -v -cover ./...

clean:
	go clean
	rm -rf bin

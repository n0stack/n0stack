# Go パラメータ
GOOS=linux
GOARCH=amd64
GOCMD=go
GOFMT=$(GOOS) fmt
GOBUILD=$(GOOS) $(GOARCH) $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


AGENT_BINARY_NAME=agent
AGGREGATER_BINARY_NAME=aggregater
API_BINARY_NAME=api
DISTRIBUTER_BINARY_NAME=distributer


all: test fmt build
dep:
	go get -u github.com/golang/dep/cmd/dep
deps:
	dep init
	dep ensure
	dep status
fmt:
	$(GOFMT)
build:
	$(GOBUILD) -o bin/$(AGENT_BINARY_NAME) -v -x 
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf bin
	rm -rf vender

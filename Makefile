binary=server
dockeruser=realwrtoff
gituser=realwrtoff
repository=go-http
version=$(shell git describe --tags)

export PATH:=${PATH}:${GOPATH}/bin
export GOPOXY=https://goproxy.io

.PHONY: all
all: vendor build test

build: cmd/*/*.go internal/*/*.go scripts/version.sh Makefile vendor
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/${binary}/main.go && \
	mkdir -p build/${repository}/bin && mv main build/${repository}/bin/${binary} && \
	mkdir -p build/${repository}/configs && cp configs/${binary}/* build/${repository}/configs && \
	mkdir -p build/${repository}/log

vendor: go.mod
	@echo "install golang dependency"
	go mod vendor

.PHONY: test
test: vendor
	@echo "Run unit tests"
	cd internal && go test -cover ./...

.PHONY: behave
behave: build
	behave features

.PHONY: clean
clean:
	rm -rf build

.PHONY: deep_clean
deep_clean:
	rm -rf build vendor
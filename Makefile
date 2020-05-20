binary=server
dockeruser=realwrtoff
gituser=realwrtoff
repository=go-http
version=$(shell git describe --tags)

export PATH:=${PATH}:${GOPATH}/bin
export GOPOXY=https://goproxy.io

.PHONY: all
all: vendor output test

output: cmd/*/*.go internal/*/*.go scripts/version.sh Makefile vendor
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/${binary}/main.go && \
	mkdir -p output/${repository}/bin && mv main output/${repository}/bin/${binary} && \
	mkdir -p output/${repository}/configs && cp configs/${binary}/* output/${repository}/configs && \
	mkdir -p output/${repository}/log

vendor: go.mod
	@echo "install golang dependency"
	go mod vendor

.PHONY: test
test: vendor
	@echo "Run unit tests"
	cd internal && go test -cover ./...

.PHONY: behave
behave: output
	behave features

.PHONY: clean
clean:
	rm -rf output

.PHONY: deep_clean
deep_clean:
	rm -rf output vendor
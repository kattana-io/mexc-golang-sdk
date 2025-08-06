GO_PATH := $(shell go env GOPATH)

lint: check-lint dep
	golangci-lint run --timeout=5m -c .golangci.yml

check-lint:
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin v1.55.2

dep:
	@go mod tidy
	@go mod download

proto:
	@docker run --platform linux/amd64 --rm -v `pwd`:/defs namely/protoc-all:1.51_1  -i ./websocket/dto/proto -d ./websocket/dto/proto -l go -o .



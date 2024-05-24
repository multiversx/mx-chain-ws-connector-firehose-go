SHELL := $(shell which bash)

debugger := $(shell which dlv)

.PHONY: help build run runb kill debug debug-ath protogen proto-install

cmd_dir = cmd/connector
binary = connector

build:
	cd ${cmd_dir} && \
		go build -v -o ${binary}

run: build
	cd ${cmd_dir} && \
		./${binary} --log-level="*:DEBUG"

debug: build
	cd ${cmd_dir} && \
		${debugger} exec ./${binary}

test:
	@echo "  >  Running unit tests"
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...

proto-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protogen:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		data/hyperOutportBlocks/hyperOutportBlock.proto data/hyperOutportBlocks/grpcBlockService.proto

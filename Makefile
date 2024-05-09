SHELL := $(shell which bash)

debugger := $(shell which dlv)

.PHONY: help build run runb kill debug debug-ath

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

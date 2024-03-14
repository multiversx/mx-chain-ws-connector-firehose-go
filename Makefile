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
		./${binary}

debug: build
	cd ${cmd_dir} && \
		${debugger} exec ./${binary}

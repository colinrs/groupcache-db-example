SHELL := /bin/bash
BASEDIR = $(shell pwd)

.PHONY: all
all: build

.PHONY: build
build: ## Build the binary file
	sh build.sh

.PHONY: run
run: ## run the binary file
	sh build.sh
	nohup ./bin/dbserver -port 9005 > ./bin/dbserver_9005.log &
	nohup ./bin/api -port 9001 > ./bin/api_9001.log &
	nohup ./bin/api -port 9002 > ./bin/api_9002.log &
	nohup ./bin/api -port 9003 > ./bin/api_9003.log &

.PHONY: help
help:
	@echo "make build - build api and dbserver"
	@echo "make run - run api and dbserver in bg"


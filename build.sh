#!/bin/bash
set -x
set -e

mkdir -p ./bin/
rm -f ./bin/dbserver
rm -f ./bin/api
go build -o  ./bin/dbserver dbserver/main.go
go build -o  ./bin/api api/main.go
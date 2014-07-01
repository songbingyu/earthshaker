#!/bin/sh

gofmt -w=true .
protoc --go_out=. ./src/netpb/*.proto
go clean
go build main.go
./main
go tool pprof --dot main main.prof > main.dot

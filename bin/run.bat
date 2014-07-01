gofmt -w=true .
protoc --go_out=. ./src/netpb/*.proto
go clean
go build main.go
main.exe --name=main
go tool pprof --dot main.exe main.prof > main.dot

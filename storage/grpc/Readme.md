## Compile proto

First you need to install go proto compiler, ensure your `go/bin` folder is included in `PATH` variable:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

then compile the proto:

```bash
cd storage/grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto
```

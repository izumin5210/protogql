module github.com/izumin5210/protogql/testdata/apis

go 1.15

require (
	github.com/izumin5210/protogql v0.0.0-00010101000000-000000000000
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/izumin5210/protogql => ../..

module apis/go/starwars

go 1.15

replace github.com/izumin5210/protogql => ../../../..

require (
	github.com/golang/protobuf v1.4.3
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)

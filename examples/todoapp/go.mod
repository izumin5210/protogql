module todoapp

go 1.15

require (
	apis/go/todo v0.0.0-00010101000000-000000000000
	apis/go/user v0.0.0-00010101000000-000000000000
	github.com/99designs/gqlgen v0.13.0
	github.com/go-task/task/v3 v3.0.0
	github.com/google/wire v0.4.0
	github.com/izumin5210/remixer v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.1.0
	google.golang.org/grpc v1.33.2
)

replace github.com/izumin5210/remixer => ../..

replace apis/go/todo => ../../testdata/apis/go/todo

replace apis/go/user => ../../testdata/apis/go/user

module github.com/izumin5210/protogql/e2e

go 1.15

require (
	apis/go/todo v0.0.0-00010101000000-000000000000
	apis/go/user v0.0.0-00010101000000-000000000000
	github.com/99designs/gqlgen v0.13.0
	github.com/google/go-cmp v0.5.4
	google.golang.org/grpc v1.33.2
	todoapp v0.0.0-00010101000000-000000000000
)

replace (
	apis/go/todo => ../testdata/apis/go/todo
	apis/go/user => ../testdata/apis/go/user
	github.com/izumin5210/protogql => ../
	todoapp => ../examples/todoapp
)

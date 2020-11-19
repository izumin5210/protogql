module todoapp

go 1.15

require (
	apis/go/todo v0.0.0-00010101000000-000000000000
	apis/go/user v0.0.0-00010101000000-000000000000
	github.com/99designs/gqlgen v0.13.0
	github.com/izumin5210/remixer v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace github.com/izumin5210/remixer => ../..

replace apis/go/todo => ../../testdata/apis/go/todo

replace apis/go/user => ../../testdata/apis/go/user

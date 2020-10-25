package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"testdata/task/graph"
	"testdata/task/graph/generated"
)

func main() {
	cfg := generated.Config{
		Resolvers: new(graph.Resolver),
	}
	cfg.Directives.Grpc = func(ctx context.Context, obj interface{}, next graphql.Resolver, service, rpc string) (interface{}, error) {
		return next(ctx)
	}
	cfg.Directives.Protobuf = func(ctx context.Context, obj interface{}, next graphql.Resolver, typeArg string) (interface{}, error) {
		return next(ctx)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		// send this panic somewhere
		log.Print(err)
		debug.PrintStack()
		return errors.New("user message on panic")
	})

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

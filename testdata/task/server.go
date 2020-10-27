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
	cfg.Directives.Grpc = func(ctx context.Context, obj interface{}, next graphql.Resolver, service string, rpc string) (res interface{}, err error) {
		return next(ctx)
	}
	cfg.Directives.Proto = func(ctx context.Context, obj interface{}, next graphql.Resolver, fullName string, packageArg string, name string, goPackage string, goName string) (res interface{}, err error) {
		return next(ctx)
	}
	cfg.Directives.ProtoField = func(ctx context.Context, obj interface{}, next graphql.Resolver, name string, typeArg string, goName string, goType string) (res interface{}, err error) {
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

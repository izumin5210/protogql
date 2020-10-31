package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/izumin5210/remixer/gqlruntime"

	"task/graph"
	"task/graph/generated"
)

func main() {
	cfg := generated.Config{
		Resolvers: new(graph.Resolver),
	}
	cfg.Directives.Grpc = gqlruntime.GrpcDirective
	cfg.Directives.Proto = gqlruntime.ProtoDirective
	cfg.Directives.ProtoField = gqlruntime.ProtoFieldDirective
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

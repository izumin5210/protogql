package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"todoapp/graph"
	"todoapp/graph/generated"
)

func main() {
	app, cleanup, err := graph.NewApp(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	cfg := generated.Config{
		Resolvers: app.Resolver,
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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("GRAPHQL_PORT"), nil))
}

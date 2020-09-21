package main

import (
	"context"
	"log"
	"os"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return protoprocessor.New(GraphQLSchemaGenerator).Process(context.Background(), os.Stdin, os.Stdout)
}

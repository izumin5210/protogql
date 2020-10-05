package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqls"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
	"github.com/vektah/gqlparser/v2/formatter"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return ProtocGenGraphQL.Run(context.Background(), os.Stdin, os.Stdout)
}

var BaseSchema = `schema {
	query: Query
	mutation: Mutation
}
type Query {
	node(id: ID!): Node
}
type Mutation {
	noop(input: NoopInput): NoopPayload
}
interface Node {
	id: ID!
}
input NoopInput {
	clientMutationId: String
}
type NoopPayload {
	clientMutationId: String
}
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @protobuf(type: String!) on OBJECT | ENUM | INPUT_OBJECT`

var ProtocGenGraphQL = &protoutil.ProtocGen{
	Generate: func(ctx context.Context, req *protoutil.ProtocGenRequest) error {
		for _, fd := range req.FilesToGenerate {
			schema, err := gqls.BuildSchema(fd)
			if err != nil {
				return err
			}

			if schema.Empty() {
				return nil
			}

			schemaDocAST, err := schema.DocumentAST()
			if err != nil {
				return err
			}

			file := req.GenerateFile(strings.TrimSuffix(fd.Path(), ".proto") + ".gql")

			f := formatter.NewFormatter(file)
			f.FormatSchemaDocument(schemaDocAST)
		}

		return nil
	},
}

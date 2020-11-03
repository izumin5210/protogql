package main

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqls"
)

func main() {
	options.Run(run)
}

var options = protogen.Options{}

func run(p *protogen.Plugin) error {
	for _, f := range p.Files {
		if !f.Generate {
			continue
		}

		schema, err := gqls.BuildSchema(f)
		if err != nil {
			return errors.Wrap(err, "failed to build GraphQL schema")
		}

		if schema.Empty() {
			return nil
		}

		schemaDocAST, err := schema.DocumentAST()
		if err != nil {
			return errors.Wrap(err, "failed to build GraphQL schema AST")
		}

		file := p.NewGeneratedFile(
			strings.TrimSuffix(f.Desc.Path(), ".proto")+".pb.graphqls",
			f.GoImportPath,
		)

		f := formatter.NewFormatter(file)
		f.FormatSchemaDocument(schemaDocAST)
	}

	return nil
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
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime`

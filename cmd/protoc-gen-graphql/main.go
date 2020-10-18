package main

import (
	"strings"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqls"
	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/compiler/protogen"
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
			return err
		}

		if schema.Empty() {
			return nil
		}

		schemaDocAST, err := schema.DocumentAST()
		if err != nil {
			return err
		}

		file := p.NewGeneratedFile(
			strings.TrimSuffix(f.Desc.Path(), ".proto")+".gql",
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
directive @protobuf(type: String!) on OBJECT | ENUM | INPUT_OBJECT`

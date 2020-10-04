package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqls"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

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

var GraphQLSchemaGenerator = protoprocessor.GenerateFunc(func(ctx context.Context, fd protoreflect.FileDescriptor) (*pluginpb.CodeGeneratorResponse_File, error) {
	schema, err := gqls.BuildSchema(fd)
	if err != nil {
		return nil, err
	}

	if schema.Empty() {
		return nil, nil
	}

	schemaDocAST, err := schema.DocumentAST()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf)
	f.FormatSchemaDocument(schemaDocAST)

	return &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(strings.TrimSuffix(fd.Path(), ".proto") + ".gql"),
		Content: proto.String(buf.String()),
	}, nil
})

package main

import (
	"bytes"
	"context"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqls"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
	"github.com/izumin5210/remixer/options"
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
	schema := &ast.SchemaDocument{}
	query := &ast.Definition{
		Kind: ast.Object,
		Name: "Query",
	}
	mutation := &ast.Definition{
		Kind: ast.Object,
		Name: "Mutation",
	}

	typeDescriptors, err := protoutil.TypeDFS(fd)
	if err != nil {
		return nil, err
	}
	gqlTypes := map[string]interface {
		gqls.Type
		gqls.Definable
	}{}
	addGQLType := func(t gqls.Type) error {
		dt, ok := t.(interface {
			gqls.Type
			gqls.Definable
		})
		if !ok {
			return nil
		}
		// TODO: should handle collisions
		gqlTypes[dt.Name()] = dt
		return nil
	}

	for _, td := range typeDescriptors {
		t, err := gqls.TypeFromProto(td)
		if err != nil {
			return nil, err
		}
		err = addGQLType(t)
		if err != nil {
			return nil, err
		}
	}

	inputTypes := []*gqls.InputObjectType{}

	err = protoutil.RangeServices(fd, func(sd protoreflect.ServiceDescriptor) error {
		err := protoutil.RangeMethods(sd, func(md protoreflect.MethodDescriptor) error {
			if q, ok := gqls.NewQuery(md); ok {
				def, err := q.FieldDefinitionAST()
				if err != nil {
					return err
				}
				query.Fields = append(query.Fields, def)
			}
			if m, ok := gqls.NewMutation(md); ok {
				def, err := m.FieldDefinitionAST()
				if err != nil {
					return err
				}
				typ, err := m.Input()
				if err != nil {
					return err
				}
				if inputType, ok := typ.(*gqls.InputObjectType); ok {
					inputTypes = append(inputTypes, inputType)
				}
				mutation.Fields = append(mutation.Fields, def)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, it := range inputTypes {
		err := addGQLType(it)
		if err != nil {
			return nil, err
		}

		itds, err := protoutil.TypeDFS(it.ProtoDescriptor())
		if err != nil {
			return nil, err
		}
		for _, itd := range itds {
			t, err := gqls.TypeFromProto(itd)
			if err != nil {
				return nil, err
			}
			if it, ok := t.(*gqls.ObjectType); ok {
				err := addGQLType(it)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	var typeNamesWillExported []string
	for n, t := range gqlTypes {
		if t.ProtoDescriptor().ParentFile() == fd {
			typeNamesWillExported = append(typeNamesWillExported, n)
		}
	}
	sort.StringSlice(typeNamesWillExported).Sort()

	for _, n := range typeNamesWillExported {
		def, err := gqlTypes[n].DefinitionAST()
		if err != nil {
			return nil, err
		}
		schema.Definitions = append(schema.Definitions, def)
	}

	if len(query.Fields) > 0 {
		schema.Extensions = append(schema.Extensions, query)
	}
	if len(mutation.Fields) > 0 {
		schema.Extensions = append(schema.Extensions, mutation)
	}

	if len(schema.Definitions) == 0 && len(schema.Extensions) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf)
	f.FormatSchemaDocument(schema)

	return &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(strings.TrimSuffix(fd.Path(), ".proto") + ".gql"),
		Content: proto.String(buf.String()),
	}, nil
})

func getQueryOptions(md protoreflect.MethodDescriptor) (*options.GraphqlQueryOptions, error) {
	ext := proto.GetExtension(md.Options(), options.E_GraphqlQuery)
	return ext.(*options.GraphqlQueryOptions), nil
}

func getMutationOptions(md protoreflect.MethodDescriptor) (*options.GraphqlMutationOptions, error) {
	ext := proto.GetExtension(md.Options(), options.E_GraphqlMutation)
	return ext.(*options.GraphqlMutationOptions), nil
}

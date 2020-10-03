package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqlschema"
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

var GraphQLSchemaGenerator = protoprocessor.GenerateFunc(func(ctx context.Context, fd protoreflect.FileDescriptor, types *protoprocessor.Types) (*pluginpb.CodeGeneratorResponse_File, error) {
	schema := &ast.SchemaDocument{}
	query := &ast.Definition{
		Kind: ast.Object,
		Name: "Query",
	}
	mutation := &ast.Definition{
		Kind: ast.Object,
		Name: "Mutation",
	}

	typeResolver := gqlschema.NewTypeResolver(types)
	typeWriter := gqlschema.NewTypeWriter(types, typeResolver)

	protoutil.RangeServices(fd, func(s protoreflect.ServiceDescriptor) error {
		protoutil.RangeMethods(s, func(m protoreflect.MethodDescriptor) error {
			qopts, err := getQueryOptions(m)
			if err != nil {
				// TODO: handing
				return err
			}
			directives := ast.DirectiveList{
				{Name: "grpc", Arguments: ast.ArgumentList{
					{Name: "service", Value: &ast.Value{Raw: string(s.FullName()), Kind: ast.StringValue}},
					{Name: "rpc", Value: &ast.Value{Raw: string(m.Name()), Kind: ast.StringValue}},
				}},
			}
			if qopts != nil {
				def := &ast.FieldDefinition{
					Name:       qopts.GetName(),
					Directives: directives,
				}

				if name := qopts.GetOutput(); name != "" {
					protoutil.RangeFields(m.Output(), func(fd protoreflect.FieldDescriptor) error {
						if string(fd.Name()) == name {
							typ, err := typeResolver.FromProto(fd)
							if err != nil {
								// TODO: handing
								return err
							}
							typeWriter.Add(typ)
							def.Type = typ.GQL
							return protoutil.BreakRange
						}
						return nil
					})
				} else {
					typ, err := typeResolver.FromMessage(m.Output())
					if err != nil {
						// TODO: handing
						return err
					}
					typeWriter.Add(typ)
					def.Type = typ.GQL
				}

				protoutil.RangeFields(m.Input(), func(fd protoreflect.FieldDescriptor) error {
					typ, err := typeResolver.FromProto(fd)
					if err != nil {
						// TODO: handing
						return err
					}
					typeWriter.Add(typ)
					def.Arguments = append(def.Arguments, typ.GQLArgumentDefinition())
					return nil
				})

				query.Fields = append(query.Fields, def)
			}
			mopts, err := getMutationOptions(m)
			if err != nil {
				// TODO: handing
				return err
			}
			if mopts != nil {
				inputType, err := typeResolver.InputFromMessage(m.Input())
				if err != nil {
					// TODO: handing
					return err
				}
				outputType, err := typeResolver.FromMessage(m.Output())
				if err != nil {
					// TODO: handing
					return err
				}

				typeWriter.AddInput(inputType)
				typeWriter.Add(outputType)

				mutation.Fields = append(mutation.Fields, &ast.FieldDefinition{
					Name: mopts.GetName(),
					Arguments: []*ast.ArgumentDefinition{
						{Name: "input", Type: inputType.GQL},
					},
					Type:       outputType.GQL,
					Directives: directives,
				})
			}
			return nil
		})

		return nil
	})

	defs, err := typeWriter.Definitions()
	if err != nil {
		// TODO: handling
		return nil, err
	}

	schema.Definitions = append(schema.Definitions, defs...)

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

package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/gqlschema"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
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

var GraphQLSchemaGenerator = protoprocessor.GenerateFunc(func(ctx context.Context, file string, types *protoprocessor.Types) (*plugin.CodeGeneratorResponse_File, error) {
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

	fd := types.FindFile(file)

	for _, s := range fd.GetService() {
		for _, m := range s.GetMethod() {
			qopts, err := getQueryOptions(m)
			if err != nil {
				// TODO: handing
				return nil, err
			}
			directives := ast.DirectiveList{
				{Name: "grpc", Arguments: ast.ArgumentList{
					{Name: "service", Value: &ast.Value{Raw: fmt.Sprintf(".%s.%s", fd.GetPackage(), s.GetName()), Kind: ast.StringValue}},
					{Name: "rpc", Value: &ast.Value{Raw: m.GetName(), Kind: ast.StringValue}},
				}},
			}
			if qopts != nil {
				def := &ast.FieldDefinition{
					Name:       qopts.GetName(),
					Directives: directives,
				}

				outputMsg := types.FindMessage(m.GetOutputType())

				if name := qopts.GetOutput(); name != "" {
					for _, fd := range outputMsg.GetField() {
						if fd.GetName() == name {
							typ, err := typeResolver.FromProto(fd)
							if err != nil {
								// TODO: handing
								return nil, err
							}
							typeWriter.Add(typ)
							def.Type = typ.GQL
							break
						}
					}
				} else {
					typ, err := typeResolver.FromProtoName(m.GetOutputType())
					if err != nil {
						// TODO: handing
						return nil, err
					}
					typeWriter.Add(typ)
					def.Type = typ.GQL
				}

				inputMsg := types.FindMessage(m.GetInputType())
				for _, fd := range inputMsg.GetField() {
					typ, err := typeResolver.FromProto(fd)
					if err != nil {
						// TODO: handing
						return nil, err
					}
					typeWriter.Add(typ)
					def.Arguments = append(def.Arguments, typ.GQLArgumentDefinition())
				}

				query.Fields = append(query.Fields, def)
			}
			mopts, err := getMutationOptions(m)
			if err != nil {
				// TODO: handing
				return nil, err
			}
			if mopts != nil {
				inputType, err := typeResolver.InputFromProtoName(m.GetInputType())
				if err != nil {
					// TODO: handing
					return nil, err
				}
				outputType, err := typeResolver.FromProtoName(m.GetOutputType())
				if err != nil {
					// TODO: handing
					return nil, err
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
		}
	}

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

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(strings.TrimSuffix(file, ".proto") + ".gql"),
		Content: proto.String(buf.String()),
	}, nil
})

func getQueryOptions(md *descriptor.MethodDescriptorProto) (*options.GraphqlQueryOptions, error) {
	ext, err := proto.GetExtension(md.GetOptions(), options.E_GraphqlQuery)
	if err == proto.ErrMissingExtension {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ext.(*options.GraphqlQueryOptions), nil
}

func getMutationOptions(md *descriptor.MethodDescriptorProto) (*options.GraphqlMutationOptions, error) {
	ext, err := proto.GetExtension(md.GetOptions(), options.E_GraphqlMutation)
	if err == proto.ErrMissingExtension {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ext.(*options.GraphqlMutationOptions), nil
}

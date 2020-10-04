package gqls

import (
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/options"
)

type Mutation struct {
	Proto   protoreflect.MethodDescriptor
	Options *options.GraphqlMutationOptions
}

func NewMutation(md protoreflect.MethodDescriptor) (*Mutation, bool) {
	if !proto.HasExtension(md.Options(), options.E_GraphqlMutation) {
		return nil, false
	}
	ext := proto.GetExtension(md.Options(), options.E_GraphqlMutation)
	opts, ok := ext.(*options.GraphqlMutationOptions)
	if !ok {
		return nil, false
	}
	return &Mutation{Proto: md, Options: opts}, true
}

func (m *Mutation) FieldDefinitionAST() (*ast.FieldDefinition, error) {
	def := &ast.FieldDefinition{
		Name:       m.Options.GetName(),
		Directives: rpcDirectivesAST(m.Proto),
	}

	inputType, err := m.Input()
	if err != nil {
		return nil, err
	}
	def.Arguments = ast.ArgumentDefinitionList{
		{
			Name: "input",
			Type: inputType.TypeAST(),
		},
	}
	outputType, err := TypeFromProto(m.Proto.Output())
	if err != nil {
		return nil, err
	}
	def.Type = outputType.TypeAST()

	return def, nil
}

func (m *Mutation) Input() (Type, error) {
	typ, err := TypeFromProto(m.Proto.Input())
	if err != nil {
		return nil, err
	}
	if objType, ok := typ.(*ObjectType); ok {
		return NewInputObjectType(objType), nil
	}
	return typ, nil
}

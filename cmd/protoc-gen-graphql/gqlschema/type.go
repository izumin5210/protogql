package gqlschema

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	scalarTypes = map[string]struct{}{
		gqlFloatType().Name():   {},
		gqlIntType().Name():     {},
		gqlStringType().Name():  {},
		gqlBooleanType().Name(): {},
		gqlIDType().Name():      {},
	}
)

func gqlTypeToInput(t *ast.Type) *ast.Type {
	if t.NamedType == "" {
		t.Elem = gqlTypeToInput(t.Elem)
		return t
	}
	if strings.HasSuffix(t.NamedType, "Request") {
		t.NamedType = strings.TrimSuffix(t.NamedType, "Request")
	}
	t.NamedType += "Input"
	return t
}

func typeFromFieldDescriptor(gqlType *ast.Type, fd *descriptor.FieldDescriptorProto) *Type {
	return &Type{GQL: gqlType, Proto: &ProtoType{Name: fd.GetTypeName(), FieldDescriptor: fd}}
}

func typeFromMessageTypeName(gqlType *ast.Type, msgTypeName string) *Type {
	return &Type{GQL: gqlType, Proto: &ProtoType{Name: msgTypeName}}
}

type Type struct {
	GQL   *ast.Type
	Proto *ProtoType
}

func (t *Type) IsScalar() bool {
	_, ok := scalarTypes[t.GQL.Name()]
	return ok
}

func (t *Type) Input() *Type {
	if t.IsScalar() {
		return t
	}
	t.GQL = gqlTypeToInput(t.GQL)
	return t
}

type ProtoType struct {
	Name            string
	FieldDescriptor *descriptor.FieldDescriptorProto
}

func (t *Type) GQLDirectives() ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "protobuf", Arguments: ast.ArgumentList{
			{Name: "type", Value: &ast.Value{Raw: t.Proto.Name, Kind: ast.StringValue}},
		}},
	}
}

func (t *Type) GQLArgumentDefinition() *ast.ArgumentDefinition {
	return &ast.ArgumentDefinition{
		Name: strcase.ToLowerCamel(t.Proto.FieldDescriptor.GetName()),
		Type: t.GQL,
	}
}

func (t *Type) GQLFieldDefinition() *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Name: strcase.ToLowerCamel(t.Proto.FieldDescriptor.GetName()),
		Type: t.GQL,
	}
}

package gqls

import (
	"github.com/izumin5210/protogql/codegen/protoutil"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Type      = (*EnumType)(nil)
	_ Definable = (*EnumType)(nil)
	_ ProtoType = (*EnumType)(nil)
)

type EnumType struct {
	Proto *protogen.Enum
}

func NewEnumType(proto *protogen.Enum) *EnumType             { return &EnumType{Proto: proto} }
func (t *EnumType) Name() string                             { return nameWithParent(t.Proto.Desc) }
func (t *EnumType) IsNullable() bool                         { return false }
func (t *EnumType) IsList() bool                             { return false }
func (t *EnumType) TypeAST() *ast.Type                       { return ast.NonNullNamedType(t.Name(), nil) }
func (t *EnumType) ProtoDescriptor() protoreflect.Descriptor { return t.Proto.Desc }
func (t *EnumType) GoIdent() protogen.GoIdent                { return t.Proto.GoIdent }

func (t *EnumType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:        ast.Enum,
		Name:        string(t.Name()),
		Directives:  enumDirectivesAST(t.Proto),
		Description: protoutil.FormatComments(t.Proto.Comments),
	}

	for _, ev := range t.Proto.Values {
		def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
			Name:        string(ev.Desc.Name()),
			Description: protoutil.FormatComments(ev.Comments),
		})
	}

	return def, nil
}

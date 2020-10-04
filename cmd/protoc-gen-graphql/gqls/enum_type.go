package gqls

import (
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Type      = (*EnumType)(nil)
	_ Definable = (*EnumType)(nil)
)

type EnumType struct {
	Proto protoreflect.EnumDescriptor
}

func newEnumType(proto protoreflect.EnumDescriptor) *EnumType { return &EnumType{Proto: proto} }
func (t *EnumType) Name() string                              { return nameWithParent(t.Proto) }
func (t *EnumType) IsNullable() bool                          { return false }
func (t *EnumType) IsList() bool                              { return false }
func (t *EnumType) TypeAST() *ast.Type                        { return ast.NonNullNamedType(t.Name(), nil) }
func (t *EnumType) ProtoDescriptor() protoreflect.Descriptor  { return t.Proto }

func (t *EnumType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:       ast.Enum,
		Name:       string(t.Name()),
		Directives: definitionDelectivesAST(t.Proto),
	}

	err := protoutil.RangeEnumValues(t.Proto, func(evd protoreflect.EnumValueDescriptor) error {
		def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
			Name: string(evd.Name()),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return def, nil
}

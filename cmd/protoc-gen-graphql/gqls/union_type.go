package gqls

import (
	"github.com/iancoleman/strcase"
	"github.com/izumin5210/remixer/codegen/protoutil"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Type      = (*UnionType)(nil)
	_ Definable = (*UnionType)(nil)
	_ ProtoType = (*UnionType)(nil)
)

type UnionType struct {
	Proto *protogen.Oneof
}

func NewUnionType(proto *protogen.Oneof) *UnionType { return &UnionType{Proto: proto} }
func (t *UnionType) Name() string {
	return nameWithParent(t.Proto.Parent.Desc) + strcase.ToCamel(string(t.Proto.Desc.Name()))
}
func (t *UnionType) IsNullable() bool                         { return false }
func (t *UnionType) IsList() bool                             { return false }
func (t *UnionType) TypeAST() *ast.Type                       { return ast.NonNullNamedType(t.Name(), nil) }
func (t *UnionType) ProtoDescriptor() protoreflect.Descriptor { return t.Proto.Desc }
func (t *UnionType) GoIdent() protogen.GoIdent                { return t.Proto.GoIdent }

func (t *UnionType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:        ast.Union,
		Name:        string(t.Name()),
		Description: protoutil.FormatComments(t.Proto.Comments),
	}

	for _, f := range t.Proto.Fields {
		ft, err := TypeFromProtoField(f)
		if err != nil {
			return nil, err
		}

		def.Types = append(def.Types, ft.Name())
	}

	def.Directives = oneofDirectivesAST(t.Proto, def.Types)

	return def, nil
}

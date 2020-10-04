package gqls

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Type      = (*InputObjectType)(nil)
	_ Definable = (*InputObjectType)(nil)
)

func NewInputObjectType(objType *ObjectType) *InputObjectType {
	return &InputObjectType{base: objType}
}

type InputObjectType struct {
	base *ObjectType
}

func (t *InputObjectType) Name() string {
	n := t.base.Name()
	n = strings.TrimSuffix(n, "Request")
	return n + "Input"
}

func (t *InputObjectType) IsNullable() bool                         { return t.base.IsNullable() }
func (t *InputObjectType) IsList() bool                             { return t.base.IsList() }
func (t *InputObjectType) TypeAST() *ast.Type                       { return ast.NonNullNamedType(t.Name(), nil) }
func (t *InputObjectType) ProtoDescriptor() protoreflect.Descriptor { return t.base.Proto }

func (t *InputObjectType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:       ast.InputObject,
		Name:       string(t.Name()),
		Directives: definitionDelectivesAST(t.base.Proto),
	}

	err := protoutil.RangeFields(t.base.Proto, func(fd protoreflect.FieldDescriptor) error {
		ft, err := TypeFromProtoField(fd)
		if err != nil {
			return err
		}
		if ot, ok := ft.(*ObjectType); ok {
			ft = NewInputObjectType(ot)
		}
		def.Fields = append(def.Fields, &ast.FieldDefinition{
			Name: strcase.ToLowerCamel(string(fd.Name())),
			Type: ft.TypeAST(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return def, nil
}

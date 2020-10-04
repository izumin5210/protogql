package gqls

import (
	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
)

var (
	_ Type      = (*ObjectType)(nil)
	_ Definable = (*ObjectType)(nil)
)

type ObjectType struct {
	Proto protoreflect.MessageDescriptor
}

func newObjectType(proto protoreflect.MessageDescriptor) *ObjectType { return &ObjectType{Proto: proto} }
func (t *ObjectType) Name() string                                   { return nameWithParent(t.Proto) }
func (t *ObjectType) IsNullable() bool                               { return false }
func (t *ObjectType) IsList() bool                                   { return false }
func (t *ObjectType) TypeAST() *ast.Type                             { return ast.NonNullNamedType(t.Name(), nil) }
func (t *ObjectType) ProtoDescriptor() protoreflect.Descriptor       { return t.Proto }

func (t *ObjectType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:       ast.Object,
		Name:       string(t.Name()),
		Directives: definitionDelectivesAST(t.Proto),
	}

	err := protoutil.RangeFields(t.Proto, func(fd protoreflect.FieldDescriptor) error {
		ft, err := TypeFromProtoField(fd)
		if err != nil {
			return err
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

package gqls

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/izumin5210/remixer/codegen/protoutil"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Type      = (*InputObjectType)(nil)
	_ Definable = (*InputObjectType)(nil)
	_ ProtoType = (*InputObjectType)(nil)
)

func NewInputObjectType(objType *ObjectType) *InputObjectType {
	return &InputObjectType{base: objType}
}

type InputObjectType struct {
	base *ObjectType
}

func (t *InputObjectType) Name() string {
	n := t.base.Name()
	return n + "Input"
}

func (t *InputObjectType) IsNullable() bool                         { return t.base.IsNullable() }
func (t *InputObjectType) IsList() bool                             { return t.base.IsList() }
func (t *InputObjectType) TypeAST() *ast.Type                       { return ast.NonNullNamedType(t.Name(), nil) }
func (t *InputObjectType) ProtoDescriptor() protoreflect.Descriptor { return t.base.Proto.Desc }
func (t *InputObjectType) GoIdent() protogen.GoIdent                { return t.base.Proto.GoIdent }

func (t *InputObjectType) DefinitionAST() (*ast.Definition, error) {
	def := &ast.Definition{
		Kind:        ast.InputObject,
		Name:        string(t.Name()),
		Directives:  messageDirectivesAST(t.base.Proto),
		Description: protoutil.FormatComments(t.base.Proto.Comments),
	}

	getFieldName := func(f *protogen.Field) string {
		return strcase.ToLowerCamel(string(f.Desc.Name()))
	}

	for _, f := range t.base.Proto.Fields {
		ft, err := TypeFromProtoField(f)
		if err != nil {
			return nil, err
		}
		origType := ft
		for {
			if mt, ok := origType.(ModifiedType); ok {
				origType = mt.Original()
			} else {
				break
			}
		}
		if ot, ok := origType.(*ObjectType); ok {
			origType = NewInputObjectType(ot)
		}
		if ft.IsNullable() || f.Oneof != nil {
			origType = NullableType(origType)
		}
		if ft.IsList() {
			origType = ListType(origType)
		}
		ft = origType
		desc := protoutil.FormatComments(f.Comments)
		if f.Oneof != nil {
			if desc != "" {
				desc += "\n\n"
			}
			fields := make([]string, 0, len(f.Oneof.Fields))
			for _, of := range f.Oneof.Fields {
				fields = append(fields, fmt.Sprintf("`%s`", getFieldName(of)))
			}
			desc += fmt.Sprintf(
				"At most one field, %s, and %s will be set at the same time.",
				strings.Join(fields[0:len(fields)-1], ", "), fields[len(fields)-1],
			)
		}
		def.Fields = append(def.Fields, &ast.FieldDefinition{
			Name:        getFieldName(f),
			Type:        ft.TypeAST(),
			Directives:  inputFieldDirectivesAST(f, ft),
			Description: desc,
		})
	}

	return def, nil
}

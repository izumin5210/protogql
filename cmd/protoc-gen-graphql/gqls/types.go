package gqls

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Type interface {
	Name() string
	IsNullable() bool
	IsList() bool
	TypeAST() *ast.Type
}

type Definable interface {
	DefinitionAST() (*ast.Definition, error)
	ProtoDescriptor() protoreflect.Descriptor
}

type ModifiedType interface {
	Original() Type
}

var (
	_ ModifiedType = (*nullableType)(nil)
	_ ModifiedType = (*listType)(nil)
)

func TypeFromProtoField(f *protogen.Field) (Type, error) {
	typ, err := rawTypeFromProtoField(f)
	if err != nil {
		return nil, err
	}
	if f.Desc.IsList() {
		typ = ListType(typ)
	}
	if f.Desc.HasOptionalKeyword() {
		typ = NullableType(typ)
	}
	return typ, nil
}

func rawTypeFromProtoField(f *protogen.Field) (Type, error) {
	switch f.Desc.Kind() {
	case protoreflect.MessageKind:
		return TypeFromProtoMessage(f.Message)
	case protoreflect.EnumKind:
		return NewEnumType(f.Enum), nil
	default:
		typ, ok := scalarTypeMap[f.Desc.Kind()]
		if !ok {
			return nil, fmt.Errorf("unsupported kind: %s", f.Desc.Kind())
		}
		return typ, nil
	}
}

func TypeFromProtoMessage(m *protogen.Message) (Type, error) {
	switch m.Desc.FullName() {
	case protoreflect.FullName("google.protobuf").Append("Empty"):
		return NullableType(BooleanType), nil
	}
	// TODO: handle other well-known types
	return newObjectType(m), nil
}

type nullableType struct {
	Type
}

func NullableType(el Type) Type          { return &nullableType{Type: el} }
func (t *nullableType) IsNullable() bool { return true }
func (t *nullableType) TypeAST() *ast.Type {
	typ := t.Type.TypeAST()
	typ.NonNull = false
	return typ
}
func (t *nullableType) Original() Type { return t.Type }

type listType struct {
	Type
}

func ListType(el Type) Type            { return &listType{Type: el} }
func (t *listType) IsList() bool       { return true }
func (t *listType) TypeAST() *ast.Type { return ast.NonNullListType(t.Type.TypeAST(), nil) }
func (t *listType) Original() Type     { return t.Type }

func nameWithParent(d protoreflect.Descriptor) string {
	name := string(d.Name())
	for parent := d.Parent(); parent != nil; parent = parent.Parent() {
		if _, ok := parent.(protoreflect.MessageDescriptor); !ok {
			break
		}
		name = string(parent.Name()) + name
	}
	return name
}

func UnwrapType(t Type) Type {
	for {
		if modified, ok := t.(ModifiedType); ok {
			t = modified.Original()
		} else {
			return t
		}
	}
}

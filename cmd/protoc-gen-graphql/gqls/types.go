package gqls

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
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

func TypeFromProto(d protoreflect.Descriptor) (Type, error) {
	switch d := d.(type) {
	case protoreflect.MessageDescriptor:
		return typeFromProtoMessage(d)
	case protoreflect.EnumDescriptor:
		return newEnumType(d), nil
	case protoreflect.FieldDescriptor:
		return TypeFromProtoField(d)
	default:
		return nil, fmt.Errorf("unsupported descriptor: %s", d.FullName())
	}
}

func TypeFromProtoField(fd protoreflect.FieldDescriptor) (Type, error) {
	typ, err := rawTypeFromProtoField(fd)
	if err != nil {
		return nil, err
	}
	if fd.IsList() {
		typ = ListType(typ)
	}
	if fd.HasOptionalKeyword() {
		typ = NullableType(typ)
	}
	return typ, nil
}

func rawTypeFromProtoField(fd protoreflect.FieldDescriptor) (Type, error) {
	switch fd.Kind() {
	case protoreflect.MessageKind:
		return typeFromProtoMessage(fd.Message())
	case protoreflect.EnumKind:
		return newEnumType(fd.Enum()), nil
	default:
		typ, ok := scalarTypeMap[protoutil.JSONKindFrom(fd.Kind())]
		if !ok {
			return nil, fmt.Errorf("unsupported kind: %s", fd.Kind())
		}
		return typ, nil
	}
}

func typeFromProtoMessage(md protoreflect.MessageDescriptor) (Type, error) {
	switch md.FullName() {
	case protoreflect.FullName("google.protobuf").Append("Empty"):
		return NullableType(BooleanType), nil
	}
	// TODO: handle other well-known types
	return newObjectType(md), nil
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

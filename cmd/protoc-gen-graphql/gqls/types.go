package gqls

import (
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	protogql_pb "github.com/izumin5210/protogql/protobuf/protogql"
)

type Type interface {
	Name() string
	IsNullable() bool
	IsList() bool
	TypeAST() *ast.Type
}

type Definable interface {
	ProtoType
	DefinitionAST() (*ast.Definition, error)
}

type ModifiedType interface {
	Original() Type
}

type ProtoType interface {
	GoIdent() protogen.GoIdent
	ProtoDescriptor() protoreflect.Descriptor
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
			return nil, errors.Errorf("kind %s is not supported", f.Desc.Kind())
		}
		return typ, nil
	}
}

func TypeFromProtoMessage(m *protogen.Message) (Type, error) {
	switch m.Desc.FullName() {
	case protoreflect.FullName("google.protobuf.Empty"):
		return NullableType(BooleanType), nil
	case protoreflect.FullName("google.protobuf.Int32Value"), protoreflect.FullName("google.protobuf.Int64Value"),
		protoreflect.FullName("google.protobuf.UInt32Value"), protoreflect.FullName("google.protobuf.UInt64Value"):
		return NullableType(&WrappedScalarType{ScalarType: IntType, Proto: m}), nil
	case protoreflect.FullName("google.protobuf.FloatValue"), protoreflect.FullName("google.protobuf.DoubleValue"):
		return NullableType(&WrappedScalarType{ScalarType: FloatType, Proto: m}), nil
	case protoreflect.FullName("google.protobuf.BoolValue"):
		return NullableType(&WrappedScalarType{ScalarType: BooleanType, Proto: m}), nil
	case protoreflect.FullName("google.protobuf.StringValue"):
		return NullableType(&WrappedScalarType{ScalarType: StringType, Proto: m}), nil
	case protoreflect.FullName("google.protobuf.Timestamp"):
		return &WrappedScalarType{ScalarType: DateTimeType, Proto: m}, nil
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
		switch pd := parent.(type) {
		case protoreflect.MessageDescriptor:
			name = string(pd.Name()) + name
		case protoreflect.FileDescriptor:
			ext, ok := proto.GetExtension(pd.Options(), protogql_pb.E_Schema).(*protogql_pb.GraphqlSchemaOptions)
			if !ok {
				break
			}
			if prefix := ext.GetTypePrefix(); prefix != "" {
				name = prefix + name
			}
		default:
			break
		}
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

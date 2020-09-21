package gqlschema

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

var (
	// default scalar types
	gqlFloatType   = func() *ast.Type { return ast.NonNullNamedType("Float", nil) }
	gqlIntType     = func() *ast.Type { return ast.NonNullNamedType("Int", nil) }
	gqlStringType  = func() *ast.Type { return ast.NonNullNamedType("String", nil) }
	gqlBooleanType = func() *ast.Type { return ast.NonNullNamedType("Boolean", nil) }
	gqlIDType      = func() *ast.Type { return ast.NonNullNamedType("ID", nil) }

	// custom scalar types
	gqlVoidType = func() *ast.Type { return ast.NamedType("Boolean", nil) }

	scalarTypes = map[string]struct{}{
		gqlFloatType().Name():   {},
		gqlIntType().Name():     {},
		gqlStringType().Name():  {},
		gqlBooleanType().Name(): {},
		gqlIDType().Name():      {},
	}
)

func NewTypeResolver(types *protoprocessor.Types) *TypeResolver {
	return &TypeResolver{types: types}
}

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

type TypeResolver struct {
	types *protoprocessor.Types
}

func (r *TypeResolver) FromProto(fd *descriptor.FieldDescriptorProto) (*Type, error) {
	typ, err := r.fromProto(fd)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return &Type{GQL: typ, Proto: &ProtoType{Name: fd.GetTypeName(), FieldDescriptor: fd}}, nil
}

func (r *TypeResolver) InputFromProto(fd *descriptor.FieldDescriptorProto) (*Type, error) {
	typ, err := r.FromProto(fd)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typ.Input(), nil
}

func (r *TypeResolver) fromProto(fd *descriptor.FieldDescriptorProto) (typ *ast.Type, err error) {
	defer func() {
		if fd.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			typ = ast.NonNullListType(typ, nil)
		}
		if fd.GetProto3Optional() {
			typ.NonNull = false
		}
	}()

	switch fd.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return gqlFloatType(), nil
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return gqlFloatType(), nil
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return gqlIntType(), nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return gqlIntType(), nil
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return gqlBooleanType(), nil
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return nil, fmt.Errorf("%s is not supported", fd.GetType())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return r.fromProtoName(fd.GetTypeName())
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return gqlIntType(), nil
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		if d := r.types.FindEnum(fd.GetTypeName()); d != nil {
			return ast.NonNullNamedType(d.GetName(), nil), nil
		}
		return nil, fmt.Errorf("%s is not found", fd.GetTypeName())
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return gqlIntType(), nil
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return gqlStringType(), nil
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return gqlIntType(), nil
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return gqlStringType(), nil
	default:
		return nil, fmt.Errorf("%s is unknown", fd.GetType())
	}
}

func (r *TypeResolver) FromProtoName(msgName string) (*Type, error) {
	typ, err := r.fromProtoName(msgName)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return &Type{GQL: typ, Proto: &ProtoType{Name: msgName}}, nil
}

func (r *TypeResolver) InputFromProtoName(msgName string) (*Type, error) {
	typ, err := r.FromProtoName(msgName)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typ.Input(), nil
}

func (r *TypeResolver) fromProtoName(msgName string) (typ *ast.Type, err error) {
	switch msgName {
	case ".google.protobuf.Empty":
		return gqlVoidType(), nil
	default:
		// TODO: wrapper types
		// TODO: google.protobuf.Timestamp
		if d := r.types.FindMessage(msgName); d != nil {
			return ast.NonNullNamedType(d.GetName(), nil), nil
		}
		return nil, fmt.Errorf("%s is not found", msgName)
	}
}

package gqlschema

import (
	"fmt"

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

type Type struct {
	GQL   *ast.Type
	Proto *ProtoType
}

func (t *Type) IsScalar() bool {
	_, ok := scalarTypes[t.GQL.Name()]
	return ok
}

type ProtoType struct {
	Name            string
	FieldDescriptor *descriptor.FieldDescriptorProto
}

type TypeResolver struct {
	types *protoprocessor.Types
}

func (r *TypeResolver) FromFieldDescriptor(fd *descriptor.FieldDescriptorProto) (*Type, error) {
	typ, err := r.fromFieldDescriptor(fd)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return &Type{GQL: typ, Proto: &ProtoType{Name: fd.GetTypeName(), FieldDescriptor: fd}}, nil
}

func (r *TypeResolver) fromFieldDescriptor(fd *descriptor.FieldDescriptorProto) (typ *ast.Type, err error) {
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
		return r.fromMessageTypeName(fd.GetTypeName())
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

func (r *TypeResolver) FromMessageTypeName(msgName string) (*Type, error) {
	typ, err := r.fromMessageTypeName(msgName)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return &Type{GQL: typ, Proto: &ProtoType{Name: msgName}}, nil
}

func (r *TypeResolver) fromMessageTypeName(msgName string) (typ *ast.Type, err error) {
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

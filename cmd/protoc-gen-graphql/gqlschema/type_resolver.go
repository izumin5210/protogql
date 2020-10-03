package gqlschema

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"

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
)

func NewTypeResolver(types *protoprocessor.Types) *TypeResolver {
	return &TypeResolver{types: types}
}

type TypeResolver struct {
	types *protoprocessor.Types
}

func (r *TypeResolver) FromProto(fd protoreflect.FieldDescriptor) (*Type, error) {
	typ, err := r.fromProto(fd)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typeFromFieldDescriptor(typ, fd), nil
}

func (r *TypeResolver) InputFromProto(fd protoreflect.FieldDescriptor) (*Type, error) {
	typ, err := r.FromProto(fd)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typ.Input(), nil
}

func (r *TypeResolver) fromProto(fd protoreflect.FieldDescriptor) (typ *ast.Type, err error) {
	defer func() {
		if fd.IsList() {
			typ = ast.NonNullListType(typ, nil)
		}
		if fd.HasOptionalKeyword() {
			typ.NonNull = false
		}
	}()

	switch fd.Kind() {
	case protoreflect.DoubleKind:
		return gqlFloatType(), nil
	case protoreflect.FloatKind:
		return gqlFloatType(), nil
	case protoreflect.Int64Kind:
		return gqlStringType(), nil
	case protoreflect.Uint64Kind:
		return gqlStringType(), nil
	case protoreflect.Int32Kind:
		return gqlIntType(), nil
	case protoreflect.Fixed64Kind:
		return gqlStringType(), nil
	case protoreflect.Fixed32Kind:
		return gqlIntType(), nil
	case protoreflect.BoolKind:
		return gqlBooleanType(), nil
	case protoreflect.StringKind:
		return gqlStringType(), nil
	case protoreflect.GroupKind:
		return nil, fmt.Errorf("%s is not supported", fd.Kind())
	case protoreflect.MessageKind:
		return r.fromMessage(fd.Message())
	case protoreflect.BytesKind:
		return gqlStringType(), nil
	case protoreflect.Uint32Kind:
		return gqlIntType(), nil
	case protoreflect.EnumKind:
		return ast.NonNullNamedType(string(fd.Enum().Name()), nil), nil
	case protoreflect.Sfixed32Kind:
		return gqlIntType(), nil
	case protoreflect.Sfixed64Kind:
		return gqlStringType(), nil
	case protoreflect.Sint32Kind:
		return gqlIntType(), nil
	case protoreflect.Sint64Kind:
		return gqlStringType(), nil
	default:
		return nil, fmt.Errorf("%s is unknown", fd.Kind())
	}
}

func (r *TypeResolver) FromMessage(msg protoreflect.MessageDescriptor) (*Type, error) {
	typ, err := r.fromMessage(msg)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typeFromMessageTypeName(typ, string(msg.FullName())), nil
}

func (r *TypeResolver) InputFromMessage(msg protoreflect.MessageDescriptor) (*Type, error) {
	typ, err := r.FromMessage(msg)
	if err != nil {
		// TODO: handling
		return nil, err
	}
	return typ.Input(), nil
}

func (r *TypeResolver) fromMessage(md protoreflect.MessageDescriptor) (typ *ast.Type, err error) {
	switch md.FullName() {
	case "google.protobuf.Empty":
		return gqlVoidType(), nil
	default:
		// TODO: wrapper types
		// TODO: google.protobuf.Timestamp
		return ast.NonNullNamedType(string(md.Name()), nil), nil
	}
}

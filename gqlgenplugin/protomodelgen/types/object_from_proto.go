package types

import (
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

var (
	_ ProtoType = (*ObjectFromProto)(nil)
)

type ObjectFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (u *ObjectFromProto) IsFromProto() bool { return true }

func (o *ObjectFromProto) GQLName() string {
	return o.def.Name
}

func (o *ObjectFromProto) GoType() GoType {
	return newGoModelType(o.def.Name)
}

func (o *ObjectFromProto) Godoc() string {
	return goutil.ToComment(o.def.Description)
}

func (o *ObjectFromProto) Fields() ([]*FieldFromProto, error) {
	fields := make([]*FieldFromProto, len(o.def.Fields))

	for i, f := range o.def.Fields {
		proto, err := gqlutil.ExtractProtoFieldDirective(f.Directives)
		if err != nil {
			return nil, errors.Wrapf(err, "%s has invalid directive", f.Name)
		}
		fields[i] = &FieldFromProto{gql: f, proto: proto, object: o}
	}

	return fields, nil
}

func (o *ObjectFromProto) ProtoGoType() GoType {
	return newGoType(o.proto.GoPackage, o.proto.GoName)
}

func (o *ObjectFromProto) ImplementedInterfaces() ([]Type, error) {
	types := make([]Type, len(o.def.Interfaces))

	for i, ifName := range o.def.Interfaces {
		typ := o.registry.FindInterfaceType(ifName)
		if typ == nil {
			return nil, errors.Errorf("interface %s was not found", ifName)
		}
		types[i] = typ
	}

	return types, nil
}

type FieldFromProto struct {
	gql    *ast.FieldDefinition
	proto  *gqlutil.ProtoFieldDirective
	object *ObjectFromProto
}

func (f *FieldFromProto) GQLName() string {
	return f.gql.Name
}

func (f *FieldFromProto) GoFieldName() string {
	return templates.ToGo(f.gql.Name)
}

func (f *FieldFromProto) PbGoFieldName() string {
	return f.proto.GoName
}

func (f *FieldFromProto) Godoc() string {
	return goutil.ToComment(f.gql.Description)
}

func (f *FieldFromProto) GoFieldTypeDefinition() string {
	var b strings.Builder

	if f.isList() {
		b.WriteString("[]")
	}

	switch {
	case f.isGoBuiltinType():
		b.WriteString(f.proto.GoTypeName)
	case f.isProtoWellKnownType():
		b.WriteString("*")
		b.WriteString(templates.CurrentImports.Lookup(f.proto.GoTypePackage))
		b.WriteString(".")
		b.WriteString(f.proto.GoTypeName)
	default:
		typ := f.object.registry.FindType(f.gql.Type.Name())
		b.WriteString(typ.GoType().TypeString())
	}

	return b.String()
}

func (f *FieldFromProto) FromProtoStatement(receiver string) string {
	if f.proto == nil {
		return ""
	}

	var b strings.Builder

	switch {
	case f.isGoBuiltinType(), f.isProtoWellKnownType():
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("()")
	case f.isList():
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(FromRepeatedProtoFuncName(typ))
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("())")
	default:
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(FromProtoFuncName(typ))
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("())")
	}

	return b.String()
}

func (f *FieldFromProto) ToProtoStatement(receiver string) string {
	if f.proto == nil {
		return ""
	}

	var b strings.Builder

	switch {
	case f.isGoBuiltinType(), f.isProtoWellKnownType():
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
	case f.isList():
		b.WriteString(ToRepeatedProtoFuncName(f.ProtoType()))
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	default:
		b.WriteString(ToProtoFuncName(f.ProtoType()))
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	}

	return b.String()
}

func (f *FieldFromProto) ProtoType() ProtoType {
	return f.object.registry.FindProtoType(f.gql.Type.Name())
}

func (f *FieldFromProto) IsOneof() bool {
	return f.proto != nil && f.proto.OneofName != ""
}

func (f *FieldFromProto) OneofMembers() []*UnionMemberFromProto {
	typ := f.object.registry.FindProtoType(f.gql.Type.Name())
	u, ok := typ.(*UnionFromProto)
	if !ok {
		return nil
	}

	return u.Members()
}

func (f *FieldFromProto) IsOneofMember() bool {
	return f.proto != nil && f.proto.OneofName != "" && f.proto.OneofName != f.proto.Name
}

func (f *FieldFromProto) IsDefinedInProto() bool {
	return f.proto != nil
}

func (f *FieldFromProto) isList() bool {
	return gqlutil.IsListType(f.gql.Type)
}

func (f *FieldFromProto) isGoBuiltinType() bool {
	if f.proto != nil {
		return f.proto.IsGoBuiltinType()
	}
	return gqlutil.IsBuiltinType(f.gql.Type)
}

func (f *FieldFromProto) isProtoWellKnownType() bool {
	return f.proto != nil && f.proto.IsWellKnownType()
}

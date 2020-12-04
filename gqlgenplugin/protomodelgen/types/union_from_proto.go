package types

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

var (
	_ ProtoType = (*UnionFromProto)(nil)
	_ ProtoType = (*UnionMemberFromProto)(nil)
)

type UnionFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (u *UnionFromProto) IsFromProto() bool { return true }

func (u *UnionFromProto) GQLName() string {
	return u.def.Name
}

func (u *UnionFromProto) GoType() GoType {
	return newGoModelInterfaceType(u.def.Name)
}

func (u *UnionFromProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionFromProto) ProtoGoType() GoType { return newGoModelInterfaceType("interface{}") }

func (u *UnionFromProto) Members() []*UnionMemberFromProto {
	members := make([]*UnionMemberFromProto, len(u.proto.Oneof.Fields))
	for i, f := range u.proto.Oneof.Fields {
		members[i] = &UnionMemberFromProto{union: u, proto: f, typ: u.registry.FindProtoType(f.Name)}
	}
	return members
}

type UnionMemberFromProto struct {
	union *UnionFromProto
	proto *gqlutil.ProtoDirectiveOneofField
	typ   ProtoType
}

func (m *UnionMemberFromProto) IsFromProto() bool { return true }

func (m *UnionMemberFromProto) GQLName() string { return m.proto.GoName }

func (m *UnionMemberFromProto) GoType() GoType {
	return m.typ.GoType()
}

func (m *UnionMemberFromProto) ProtoGoType() GoType {
	return wrapGoType(m.typ.ProtoGoType(), m.proto.GoName)
}

func (m *UnionMemberFromProto) PbGoTypeFieldName() string {
	return m.proto.Name
}

func (m *UnionMemberFromProto) Type() ProtoType {
	return m.typ
}

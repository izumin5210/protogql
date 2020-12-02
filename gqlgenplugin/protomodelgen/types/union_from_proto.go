package types

import (
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type UnionFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (u *UnionFromProto) GoTypeName() string {
	return u.def.Name
}

func (u *UnionFromProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionFromProto) FuncNameFromProto() string {
	return u.GoTypeName() + "FromProto"
}

func (u *UnionFromProto) PbGoTypeName() string              { panic("unreachable") }
func (u *UnionFromProto) FuncNameFromRepeatedProto() string { panic("unreachable") }
func (u *UnionFromProto) FuncNameToProto() string           { panic("unreachable") }
func (u *UnionFromProto) FuncNameToRepeatedProto() string   { panic("unreachable") }

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

func (m *UnionMemberFromProto) GoTypeName() string {
	return m.typ.GoTypeName()
}

func (m *UnionMemberFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(m.union.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(m.proto.GoName)

	return b.String()
}

func (m *UnionMemberFromProto) PbGoTypeFieldName() string {
	return m.proto.Name
}

func (m *UnionMemberFromProto) FuncNameFromProto() string {
	return m.typ.FuncNameFromProto()
}

func (m *UnionMemberFromProto) FuncNameToProto() string {
	return m.proto.GoName + "ToProto"
}

func (m *UnionMemberFromProto) Type() ProtoType {
	return m.typ
}

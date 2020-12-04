package types

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
)

var (
	_ ProtoType = (*UnionHasProto)(nil)
)

type UnionHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (u *UnionHasProto) IsFromProto() bool { return false }

func (u *UnionHasProto) GQLName() string {
	return u.def.Name
}

func (u *UnionHasProto) GoType() GoType {
	return newGoModelInterfaceType(u.def.Name)
}

func (u *UnionHasProto) ProtoGoType() GoType {
	return newGoModelType(u.GoType().Name() + "_Proto")
}

func (u *UnionHasProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionHasProto) Members() []UnionMemberType {
	members := make([]UnionMemberType, len(u.def.Types))
	for i, t := range u.def.Types {
		switch typ := u.registry.FindType(t).(type) {
		case ProtoType:
			members[i] = &UnionMemberHasProto{typ: typ}
		case Type:
			members[i] = &PlainUnionMember{typ: typ}
		default:
			panic("unreachable")
		}
	}
	return members
}

type PlainUnionMember struct {
	typ   Type
	union Type
}

func (m *PlainUnionMember) Type() Type {
	return m.typ
}

type UnionMemberHasProto struct {
	typ   ProtoType
	union *UnionHasProto
}

func (m *UnionMemberHasProto) Type() Type {
	return m.typ
}

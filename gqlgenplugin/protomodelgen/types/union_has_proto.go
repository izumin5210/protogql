package types

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
)

type UnionHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (u *UnionHasProto) GoTypeName() string {
	return u.def.Name
}

func (u *UnionHasProto) GoWrapperTypeName() string {
	return u.GoTypeName() + "_Proto"
}

func (u *UnionHasProto) FuncNameFromProto() string {
	return u.GoTypeName() + "FromProto"
}

func (u *UnionHasProto) FuncNameFromRepeatedProto() string {
	return u.GoTypeName() + "ListFromRepeatedProto"
}

func (u *UnionHasProto) FuncNameToProto() string {
	return u.GoTypeName() + "ToProto"
}

func (u *UnionHasProto) FuncNameToRepeatedProto() string {
	return u.GoTypeName() + "ListToRepeatedProto"
}

func (u *UnionHasProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionHasProto) Members() []UnionMemberType {
	members := make([]UnionMemberType, len(u.def.Types))
	for i, t := range u.def.Types {
		switch typ := u.registry.FindType(t).(type) {
		case ProtoLikeType:
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
	typ   ProtoLikeType
	union *UnionHasProto
}

func (m *UnionMemberHasProto) Type() Type {
	return m.typ
}

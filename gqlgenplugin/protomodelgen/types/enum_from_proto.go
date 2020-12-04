package types

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

var (
	_ ProtoType = (*EnumFromProto)(nil)
)

type EnumFromProto struct {
	def   *ast.Definition
	proto *gqlutil.ProtoDirective
}

func (m *EnumFromProto) IsFromProto() bool { return true }

func (e *EnumFromProto) GQLName() string {
	return e.def.Name
}

func (e *EnumFromProto) GoType() GoType {
	return newGoModelType(e.def.Name)
}

func (e *EnumFromProto) ProtoGoType() GoType {
	return newGoInterfaceType(e.proto.GoPackage, e.proto.GoName)
}

func (e *EnumFromProto) Godoc() string {
	return goutil.ToComment(e.def.Description)
}

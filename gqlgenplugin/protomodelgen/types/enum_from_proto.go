package types

import (
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type EnumFromProto struct {
	def   *ast.Definition
	proto *gqlutil.ProtoDirective
}

func (e *EnumFromProto) GQLName() string {
	return e.def.Name
}

func (e *EnumFromProto) GoTypeName() string {
	return e.def.Name
}

func (e *EnumFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(e.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(e.proto.GoName)

	return b.String()
}

func (e *EnumFromProto) Godoc() string {
	return goutil.ToComment(e.def.Description)
}

func (e *EnumFromProto) FuncNameFromProto() string {
	return e.GoTypeName() + "FromProto"
}

func (e *EnumFromProto) FuncNameFromRepeatedProto() string {
	return e.GoTypeName() + "ListFromRepeatedProto"
}

func (e *EnumFromProto) FuncNameToProto() string {
	return e.GoTypeName() + "ToProto"
}

func (e *EnumFromProto) FuncNameToRepeatedProto() string {
	return e.GoTypeName() + "ListToRepeatedProto"
}

package types

import (
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	_ Type = (*PlainObject)(nil)
	_ Type = (*PlainInterface)(nil)
)

type PlainObject struct {
	def *ast.Definition
}

func (o *PlainObject) GQLName() string {
	return o.def.Name
}

func (o *PlainObject) GoType() GoType {
	return newGoModelType(o.def.Name)
}

type PlainInterface struct {
	def *ast.Definition
}

func (it *PlainInterface) GQLName() string {
	return it.def.Name
}

func (it *PlainInterface) GoType() GoType {
	return newGoModelInterfaceType(it.def.Name)
}

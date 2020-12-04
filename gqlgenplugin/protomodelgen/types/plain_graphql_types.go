package types

import (
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	_ Type = (*PlainObject)(nil)
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

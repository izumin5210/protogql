package types

import (
	"github.com/vektah/gqlparser/v2/ast"
)

type PlainObject struct {
	def *ast.Definition
}

func (o *PlainObject) GoTypeName() string {
	return o.def.Name
}

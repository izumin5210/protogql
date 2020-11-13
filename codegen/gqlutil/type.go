package gqlutil

import (
	"github.com/vektah/gqlparser/v2/ast"
)

func IsListType(t *ast.Type) bool {
	return t.NamedType == ""
}

func IsBuiltinType(t *ast.Type) bool {
	switch t.Name() {
	case "ID", "Int", "Float", "String", "Boolean":
		return true
	default:
		return false
	}
}

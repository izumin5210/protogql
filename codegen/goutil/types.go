package goutil

import (
	"fmt"
	"go/types"
)

func GetTypePackageName(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return GetTypePackageName(t.Elem())
	case *types.Slice:
		return GetTypePackageName(t.Elem())
	case *types.Named:
		return t.Obj().Pkg().Path()
	default:
		panic(fmt.Errorf("not supported: %s", t))
	}
}

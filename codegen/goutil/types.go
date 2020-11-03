package goutil

import (
	"go/types"

	"github.com/pkg/errors"
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
		panic(errors.Errorf("%T is not supported", t))
	}
}

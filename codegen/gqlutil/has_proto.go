package gqlutil

import (
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
)

func HasProto(def *ast.Definition, defMap map[string]*ast.Definition) (resp bool, err error) {
	return hasProto(def, defMap, map[string]bool{})
}

func hasProto(def *ast.Definition, defMap map[string]*ast.Definition, cache map[string]bool) (resp bool, err error) {
	defer func() {
		if err == nil {
			cache[def.Name] = resp
		}
	}()

	if b, ok := cache[def.Name]; ok {
		return b, nil
	}
	cache[def.Name] = false

	proto, err := ExtractProtoDirective(def.Directives)
	if err != nil {
		return false, errors.Wrapf(err, "%s has invalid directive", def.Name)
	}
	if proto != nil {
		return true, nil
	}

	switch def.Kind {
	case ast.Object, ast.InputObject:
		types := make([]string, len(def.Fields))
		for i, f := range def.Fields {
			types[i] = f.Type.Name()
		}
		return typesHasProto(types, defMap, cache)
	case ast.Scalar:
		return false, nil
	case ast.Union:
		return typesHasProto(def.Types, defMap, cache)
	case ast.Enum, ast.Interface:
		panic("not supported")
	default:
		panic("unreachable")
	}
}

func typesHasProto(types []string, defMap map[string]*ast.Definition, cache map[string]bool) (resp bool, err error) {
	for _, t := range types {
		childDef, ok := defMap[t]
		if !ok {
			continue
		}
		ok, err := hasProto(childDef, defMap, cache)
		if err != nil {
			return false, errors.WithStack(err)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

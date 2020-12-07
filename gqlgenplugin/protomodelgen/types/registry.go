package types

import (
	"fmt"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type Registry struct {
	objectsFromProto map[string]*ObjectFromProto
	objectsHasProto  map[string]*ObjectHasProto
	plainObjects     map[string]*PlainObject
	enumsFromProto   map[string]*EnumFromProto
	unionsFromProto  map[string]*UnionFromProto
	unionsHasProto   map[string]*UnionHasProto
	plainInterfaces  map[string]*PlainInterface
	data             *codegen.Data
}

func CreateRegistry(data *codegen.Data) (*Registry, error) {
	return createRegistry(data, data.Schema)
}

func CreateRegistryFromSchema(schema *ast.Schema) (*Registry, error) {
	return createRegistry(nil, schema)
}

func createRegistry(data *codegen.Data, schema *ast.Schema) (*Registry, error) {
	reg := &Registry{
		objectsFromProto: map[string]*ObjectFromProto{},
		objectsHasProto:  map[string]*ObjectHasProto{},
		plainObjects:     map[string]*PlainObject{},
		enumsFromProto:   map[string]*EnumFromProto{},
		unionsFromProto:  map[string]*UnionFromProto{},
		unionsHasProto:   map[string]*UnionHasProto{},
		plainInterfaces:  map[string]*PlainInterface{},
		data:             data,
	}
	for _, def := range schema.Types {
		if strings.HasPrefix(def.Name, "__") {
			continue
		}
		if q, m := schema.Query, schema.Mutation; (q != nil && def.Name == q.Name) || (m != nil && def.Name == m.Name) {
			continue
		}

		switch def.Kind {
		case ast.Object, ast.InputObject:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.objectsFromProto[def.Name] = &ObjectFromProto{def: def, proto: proto, registry: reg}
			} else if ok, err := gqlutil.HasProto(def, schema.Types); err == nil && ok {
				reg.objectsHasProto[def.Name] = &ObjectHasProto{def: def, registry: reg}
			} else {
				reg.plainObjects[def.Name] = &PlainObject{def: def}
			}

		case ast.Enum:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.enumsFromProto[def.Name] = &EnumFromProto{def: def, proto: proto}
			} else {
				panic("Plain GraphQL Enums is not supported yet")
			}

		case ast.Scalar:
			// no-op

		case ast.Union:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.unionsFromProto[def.Name] = &UnionFromProto{def: def, proto: proto, registry: reg}
			} else if ok, err := gqlutil.HasProto(def, schema.Types); err == nil && ok {
				reg.unionsHasProto[def.Name] = &UnionHasProto{def: def, registry: reg}
			} else {
				panic("Plain GraphQL Unions is not supported yet")
			}

		case ast.Interface:
			reg.plainInterfaces[def.Name] = &PlainInterface{def: def}

		default:
			// TODO: not implemented
			panic(fmt.Errorf("%s is not supported yet", def.Kind))
		}
	}

	return reg, nil
}

func (r *Registry) FindType(name string) Type {
	if typ := r.FindProtoType(name); typ != nil {
		return typ
	}
	if obj, ok := r.plainObjects[name]; ok {
		return obj
	}

	return nil
}

func (r *Registry) FindProtoType(name string) ProtoType {
	if obj, ok := r.objectsFromProto[name]; ok {
		return obj
	}
	if enum, ok := r.enumsFromProto[name]; ok {
		return enum
	}
	if union, ok := r.unionsFromProto[name]; ok {
		return union
	}
	if obj, ok := r.objectsHasProto[name]; ok {
		return obj
	}
	if union, ok := r.unionsHasProto[name]; ok {
		return union
	}

	return nil
}

func (r *Registry) FindInterfaceType(name string) Type {
	if itf, ok := r.plainInterfaces[name]; ok {
		return itf
	}

	return nil
}

func (r *Registry) FindObjectOrInput(def *ast.Definition) *codegen.Object {
	if def.Kind == ast.InputObject {
		return r.data.Inputs.ByName(def.Name)
	}
	return r.data.Objects.ByName(def.Name)
}

func (r *Registry) ObjectsFromProto() []*ObjectFromProto {
	objs := make([]*ObjectFromProto, 0, len(r.objectsFromProto))
	for _, o := range r.objectsFromProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GQLName() < objs[j].GQLName() })

	return objs
}

func (r *Registry) ObjectsHasProto() []*ObjectHasProto {
	// FIXME
	if r.data == nil {
		return []*ObjectHasProto{}
	}

	objs := make([]*ObjectHasProto, 0, len(r.objectsHasProto))
	for _, o := range r.objectsHasProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GQLName() < objs[j].GQLName() })

	return objs
}

func (r *Registry) EnumsFromProto() []*EnumFromProto {
	enums := make([]*EnumFromProto, 0, len(r.enumsFromProto))
	for _, e := range r.enumsFromProto {
		enums = append(enums, e)
	}

	sort.Slice(enums, func(i, j int) bool { return enums[i].GQLName() < enums[j].GQLName() })

	return enums
}

func (r *Registry) UnionsFromProto() []*UnionFromProto {
	unions := make([]*UnionFromProto, 0, len(r.unionsFromProto))
	for _, u := range r.unionsFromProto {
		unions = append(unions, u)
	}

	sort.Slice(unions, func(i, j int) bool { return unions[i].GQLName() < unions[j].GQLName() })

	return unions
}

func (r *Registry) UnionsHasProto() []*UnionHasProto {
	unions := make([]*UnionHasProto, 0, len(r.unionsHasProto))
	for _, u := range r.unionsHasProto {
		unions = append(unions, u)
	}

	sort.Slice(unions, func(i, j int) bool { return unions[i].GQLName() < unions[j].GQLName() })

	return unions
}

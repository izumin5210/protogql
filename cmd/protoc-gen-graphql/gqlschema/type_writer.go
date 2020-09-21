package gqlschema

import (
	"fmt"
	"sort"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

type TypeWriter struct {
	usedTypes map[string]*typeWithKind
	defs      map[string]*ast.Definition
	types     *protoprocessor.Types
	gqlTypes  *TypeResolver
}

type typeWithKind struct {
	*Type
	Kind ast.DefinitionKind
}

func NewTypeWriter(
	types *protoprocessor.Types,
	gqlTypes *TypeResolver,
) *TypeWriter {
	return &TypeWriter{
		usedTypes: map[string]*typeWithKind{},
		defs:      map[string]*ast.Definition{},
		types:     types,
		gqlTypes:  gqlTypes,
	}
}

func (w *TypeWriter) AddInput(typ *Type) {
	w.add(typ, ast.InputObject)
}

func (w *TypeWriter) Add(typ *Type) {
	w.add(typ, ast.Object)
}

func (w *TypeWriter) add(typ *Type, kind ast.DefinitionKind) {
	if typ.IsScalar() {
		return
	}
	if ed := w.types.FindEnum(typ.Proto.Name); ed != nil {
		kind = ast.Enum
	}
	w.usedTypes[typ.Proto.Name] = &typeWithKind{Type: typ, Kind: kind}
}

func (w *TypeWriter) Definitions() ([]*ast.Definition, error) {
	for {
		allOK := true
		for typeName, typ := range w.usedTypes {
			if _, ok := w.defs[typeName]; ok {
				continue
			}
			allOK = false

			switch typ.Kind {
			case ast.Object, ast.InputObject:
				md := w.types.FindMessage(typ.Proto.Name)
				def := &ast.Definition{Kind: typ.Kind, Name: md.GetName()}

				for _, fd := range md.GetField() {
					def.Fields = append(def.Fields, &ast.FieldDefinition{
						Name: strcase.ToLowerCamel(fd.GetName()),
						Type: typ.GQL,
					})
					subtyp, err := w.gqlTypes.FromFieldDescriptor(fd)
					if err != nil {
						// TODO: handling
						return nil, err
					}
					if typ.Kind == ast.InputObject {
						w.AddInput(subtyp)
					} else {
						w.Add(subtyp)
					}
				}

				w.defs[typ.Proto.Name] = def

			case ast.Enum:
				ed := w.types.FindEnum(typ.Proto.Name)
				def := &ast.Definition{
					Kind: ast.Enum,
					Name: ed.GetName(),
				}

				for _, evd := range ed.GetValue() {
					def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
						Name: evd.GetName(),
					})
				}

				w.defs[typ.Proto.Name] = def

			default:
				return nil, fmt.Errorf("%s(kind=%s) was cannot proceeded", typ.Proto.Name, typ.Kind)
			}
		}
		if allOK {
			break
		}
	}

	defs := make([]*ast.Definition, 0, len(w.defs))
	for _, def := range w.defs {
		defs = append(defs, def)
	}

	sort.Slice(defs, func(i, j int) bool { return defs[i].Name < defs[j].Name })

	return defs, nil
}

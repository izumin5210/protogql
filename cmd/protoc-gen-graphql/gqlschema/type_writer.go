package gqlschema

import (
	"fmt"
	"sort"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

type TypeWriter struct {
	generalTypes map[string]*Type
	inputTypes   map[string]*Type
	defs         map[string]*ast.Definition
	types        *protoprocessor.Types
	gqlTypes     *TypeResolver
}

func NewTypeWriter(
	types *protoprocessor.Types,
	gqlTypes *TypeResolver,
) *TypeWriter {
	return &TypeWriter{
		generalTypes: map[string]*Type{},
		inputTypes:   map[string]*Type{},
		defs:         map[string]*ast.Definition{},
		types:        types,
		gqlTypes:     gqlTypes,
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
	if kind == ast.InputObject {
		w.inputTypes[typ.Proto.Name] = typ
	} else {
		w.generalTypes[typ.Proto.Name] = typ
	}
}

func (w *TypeWriter) Definitions() ([]*ast.Definition, error) {
	for {
		allOK := true
		for typeName, typ := range w.generalTypes {
			if _, ok := w.defs[typeName]; ok {
				continue
			}
			allOK = false

			if md := w.types.FindMessage(typ.Proto.Name); md != nil {
				def := &ast.Definition{
					Kind: ast.Object,
					Name: md.GetName(),
					Directives: ast.DirectiveList{
						{Name: "protobuf", Arguments: ast.ArgumentList{
							{Name: "type", Value: &ast.Value{Raw: typ.Proto.Name, Kind: ast.StringValue}},
						}},
					},
				}

				for _, fd := range md.GetField() {
					subtyp, err := w.gqlTypes.FromProto(fd)
					if err != nil {
						// TODO: handling
						return nil, err
					}
					def.Fields = append(def.Fields, &ast.FieldDefinition{
						Name: strcase.ToLowerCamel(fd.GetName()),
						Type: subtyp.GQL,
					})
					w.Add(subtyp)
				}

				w.defs[typ.Proto.Name] = def
			} else if ed := w.types.FindEnum(typ.Proto.Name); ed != nil {
				def := &ast.Definition{
					Kind: ast.Enum,
					Name: ed.GetName(),
					Directives: ast.DirectiveList{
						{Name: "protobuf", Arguments: ast.ArgumentList{
							{Name: "type", Value: &ast.Value{Raw: typ.Proto.Name, Kind: ast.StringValue}},
						}},
					},
				}

				for _, evd := range ed.GetValue() {
					def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
						Name: evd.GetName(),
					})
				}

				w.defs[typ.Proto.Name] = def
			} else {
				return nil, fmt.Errorf("%s was cannot proceeded", typ.Proto.Name)
			}
		}
		for typeName, typ := range w.inputTypes {
			if _, ok := w.defs[typeName]; ok {
				continue
			}
			allOK = false
			md := w.types.FindMessage(typ.Proto.Name)
			def := &ast.Definition{
				Kind: ast.InputObject,
				Name: typ.GQL.Name(),
				Directives: ast.DirectiveList{
					{Name: "protobuf", Arguments: ast.ArgumentList{
						{Name: "type", Value: &ast.Value{Raw: typ.Proto.Name, Kind: ast.StringValue}},
					}},
				},
			}

			for _, fd := range md.GetField() {
				subtyp, err := w.gqlTypes.InputFromProto(fd)
				if err != nil {
					// TODO: handling
					return nil, err
				}
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: strcase.ToLowerCamel(fd.GetName()),
					Type: subtyp.GQL,
				})
				w.AddInput(subtyp)
			}

			w.defs[typ.Proto.Name] = def
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

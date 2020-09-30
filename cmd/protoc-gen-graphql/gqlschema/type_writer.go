package gqlschema

import (
	"fmt"
	"sort"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

type TypeWriter struct {
	objectTypes map[string]*Type
	inputTypes  map[string]*Type
	enumTypes   map[string]*Type
	defs        map[string]*ast.Definition
	types       *protoprocessor.Types
	gqlTypes    *TypeResolver
}

func NewTypeWriter(
	types *protoprocessor.Types,
	gqlTypes *TypeResolver,
) *TypeWriter {
	return &TypeWriter{
		objectTypes: map[string]*Type{},
		inputTypes:  map[string]*Type{},
		enumTypes:   map[string]*Type{},
		defs:        map[string]*ast.Definition{},
		types:       types,
		gqlTypes:    gqlTypes,
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
	switch kind {
	case ast.Object:
		w.objectTypes[typ.Proto.Name] = typ
	case ast.InputObject:
		w.inputTypes[typ.Proto.Name] = typ
	case ast.Enum:
		w.enumTypes[typ.Proto.Name] = typ
	default:
		panic("unreachable")
	}
}

func (w *TypeWriter) Definitions() ([]*ast.Definition, error) {
	var err error

	err = w.buildDefinitionsRecursively(w.objectTypes, w.buildObjectDefinition)
	if err != nil {
		return nil, err
	}
	err = w.buildDefinitionsRecursively(w.inputTypes, w.buildInputObjectDefinition)
	if err != nil {
		return nil, err
	}
	err = w.buildDefinitionsRecursively(w.enumTypes, w.buildEnumDefinition)
	if err != nil {
		return nil, err
	}

	defs := make([]*ast.Definition, 0, len(w.defs))
	for _, def := range w.defs {
		defs = append(defs, def)
	}

	sort.Slice(defs, func(i, j int) bool { return defs[i].Name < defs[j].Name })

	return defs, nil
}

func (w *TypeWriter) buildDefinitionsRecursively(types map[string]*Type, build func(*Type) (*ast.Definition, error)) error {
	for {
		allOK := true
		for typeName, typ := range types {
			if _, ok := w.defs[typeName]; ok {
				continue
			}
			allOK = false

			def, err := build(typ)
			if err != nil {
				return err
			}

			w.defs[typ.Proto.Name] = def
		}
		if allOK {
			return nil
		}
	}
}

func (w *TypeWriter) buildEnumDefinition(typ *Type) (*ast.Definition, error) {
	ed := w.types.FindEnum(typ.Proto.Name)
	if ed == nil {
		return nil, fmt.Errorf("enum %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.Enum,
		Name:       ed.GetName(),
		Directives: typ.GQLDirectives(),
	}

	for _, evd := range ed.GetValue() {
		def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
			Name: evd.GetName(),
		})
	}

	return def, nil
}

func (w *TypeWriter) buildObjectDefinition(typ *Type) (*ast.Definition, error) {
	md := w.types.FindMessage(typ.Proto.Name)
	if md == nil {
		return nil, fmt.Errorf("message %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.Object,
		Name:       md.GetName(),
		Directives: typ.GQLDirectives(),
	}

	for _, fd := range md.GetField() {
		subtyp, err := w.gqlTypes.FromProto(fd)
		if err != nil {
			// TODO: handling
			return nil, err
		}
		def.Fields = append(def.Fields, subtyp.GQLFieldDefinition())
		w.Add(subtyp)
	}

	return def, nil
}

func (w *TypeWriter) buildInputObjectDefinition(typ *Type) (*ast.Definition, error) {
	md := w.types.FindMessage(typ.Proto.Name)
	if md == nil {
		return nil, fmt.Errorf("message %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.InputObject,
		Name:       typ.GQL.Name(),
		Directives: typ.GQLDirectives(),
	}

	for _, fd := range md.GetField() {
		subtyp, err := w.gqlTypes.InputFromProto(fd)
		if err != nil {
			// TODO: handling
			return nil, err
		}
		def.Fields = append(def.Fields, subtyp.GQLFieldDefinition())
		w.AddInput(subtyp)
	}

	return def, nil
}

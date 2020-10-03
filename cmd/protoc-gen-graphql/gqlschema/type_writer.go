package gqlschema

import (
	"fmt"
	"sort"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"

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
	if _, err := w.types.FindEnumByName(protoreflect.FullName(typ.Proto.Name)); err == nil {
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
	ed, err := w.types.FindEnumByName(protoreflect.FullName(typ.Proto.Name))
	if err != nil {
		return nil, fmt.Errorf("enum %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.Enum,
		Name:       string(ed.Name()),
		Directives: typ.GQLDirectives(),
	}

	values := ed.Values()
	n := values.Len()
	for i := 0; i < n; i++ {
		evd := values.Get(i)
		def.EnumValues = append(def.EnumValues, &ast.EnumValueDefinition{
			Name: string(evd.Name()),
		})
	}

	return def, nil
}

func (w *TypeWriter) buildObjectDefinition(typ *Type) (*ast.Definition, error) {
	md, err := w.types.FindMessageByName(protoreflect.FullName(typ.Proto.Name))
	if err != nil {
		return nil, fmt.Errorf("message %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.Object,
		Name:       string(md.Name()),
		Directives: typ.GQLDirectives(),
	}

	fields := md.Fields()
	n := fields.Len()
	for i := 0; i < n; i++ {
		fd := fields.Get(i)
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
	md, err := w.types.FindMessageByName(protoreflect.FullName(typ.Proto.Name))
	if err != nil {
		return nil, fmt.Errorf("message %s is not found", typ.Proto.Name)
	}

	def := &ast.Definition{
		Kind:       ast.InputObject,
		Name:       typ.GQL.Name(),
		Directives: typ.GQLDirectives(),
	}

	fields := md.Fields()
	n := fields.Len()
	for i := 0; i < n; i++ {
		fd := fields.Get(i)
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

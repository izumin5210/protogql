package gqls

import (
	"sort"

	"github.com/vektah/gqlparser/v2/ast"
)

type SchemaAST struct {
	*Schema
}

func NewSchemaAST(s *Schema) *SchemaAST {
	return &SchemaAST{
		Schema: s,
	}
}

func (s *SchemaAST) Build() (*ast.SchemaDocument, error) {
	doc := &ast.SchemaDocument{}

	defs, err := s.typeDefinitionsAST()
	if err != nil {
		return nil, err
	}
	doc.Definitions = append(doc.Definitions, defs...)

	return doc, nil
}

func (s *SchemaAST) typeDefinitionsAST() ([]*ast.Definition, error) {
	var defs []*ast.Definition

	for _, t := range s.Types {
		def, err := t.DefinitionAST()
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	sort.Slice(defs, func(i, j int) bool {
		return defs[i].Name < defs[j].Name
	})

	return defs, nil
}

func (s *SchemaAST) queriesDefinitionAST() (*ast.Definition, error) {
	if len(s.Queries) == 0 {
		return nil, nil
	}
	def := &ast.Definition{
		Kind: ast.Object,
		Name: "Query",
	}
	for _, q := range s.Queries {
		f, err := q.FieldDefinitionAST()
		if err != nil {
			return nil, err
		}
		def.Fields = append(def.Fields, f)
	}
	return def, nil
}

func (s *SchemaAST) mutationsDefinitionAST() (*ast.Definition, error) {
	if len(s.Mutations) == 0 {
		return nil, nil
	}
	def := &ast.Definition{
		Kind: ast.Object,
		Name: "Mutation",
	}
	for _, m := range s.Mutations {
		f, err := m.FieldDefinitionAST()
		if err != nil {
			return nil, err
		}
		def.Fields = append(def.Fields, f)
	}
	return def, nil
}

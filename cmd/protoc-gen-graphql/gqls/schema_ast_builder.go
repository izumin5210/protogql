package gqls

import (
	"sort"

	"github.com/pkg/errors"
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
		return nil, errors.WithStack(err)
	}
	doc.Definitions = append(doc.Definitions, defs...)

	return doc, nil
}

func (s *SchemaAST) typeDefinitionsAST() ([]*ast.Definition, error) {
	var defs []*ast.Definition

	for _, t := range s.Types {
		def, err := t.DefinitionAST()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build %s definition AST", t.Name())
		}
		defs = append(defs, def)
	}
	sort.Slice(defs, func(i, j int) bool {
		return defs[i].Name < defs[j].Name
	})

	return defs, nil
}

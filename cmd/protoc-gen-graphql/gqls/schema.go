package gqls

import (
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildSchema(fd protoreflect.FileDescriptor) (*Schema, error) {
	return NewSchemaBuilder().Build(fd)
}

func newSchema() *Schema {
	return &Schema{
		Types: map[string]interface {
			Type
			Definable
		}{}}
}

type Schema struct {
	Queries   []*Query
	Mutations []*Mutation
	Types     map[string]interface {
		Type
		Definable
	}
}

func (s *Schema) Empty() bool {
	return len(s.Types) == 0 && len(s.Queries) == 0 && len(s.Mutations) == 0
}

func (s *Schema) DocumentAST() (*ast.SchemaDocument, error) {
	return NewSchemaAST(s).Build()
}

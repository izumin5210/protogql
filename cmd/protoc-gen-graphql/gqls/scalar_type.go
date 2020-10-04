package gqls

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
)

var (
	FloatType   = newScalarType("Float")
	IntType     = newScalarType("Int")
	StringType  = newScalarType("String")
	BooleanType = newScalarType("Boolean")
	IDType      = newScalarType("ID")

	scalarTypeMap = map[protoutil.JSONKind]Type{
		protoutil.JSONInt:          IntType,
		protoutil.JSONFloat:        FloatType,
		protoutil.JSONString:       StringType,
		protoutil.JSONBoolean:      BooleanType,
		protoutil.JSONBase64String: StringType,
	}
)

type scalarType struct {
	name string
}

func newScalarType(name string) Type     { return &scalarType{name: name} }
func (t *scalarType) Name() string       { return t.name }
func (t *scalarType) IsNullable() bool   { return false }
func (t *scalarType) IsList() bool       { return false }
func (t *scalarType) TypeAST() *ast.Type { return ast.NonNullNamedType(t.Name(), nil) }

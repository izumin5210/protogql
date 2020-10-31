package gqls

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/protoutil"
)

var (
	FloatType    = newScalarType("Float")
	IntType      = newScalarType("Int")
	StringType   = newScalarType("String")
	BooleanType  = newScalarType("Boolean")
	IDType       = newScalarType("ID")
	DateTimeType = newScalarType("DateTime")

	scalarTypeMap = map[protoreflect.Kind]Type{}

	_ ProtoType = (*WrappedScalarType)(nil)
)

func init() {
	for protoKind, jsonKind := range protoutil.JSONKindMap {
		var name string
		switch jsonKind {
		case protoutil.JSONFloat:
			name = "Float"
		case protoutil.JSONInt:
			name = "Int"
		case protoutil.JSONString:
			name = "String"
		case protoutil.JSONBoolean:
			name = "Boolean"
		case protoutil.JSONBase64String:
			name = "String"
		}

		if name == "String" && strings.HasSuffix(protoKind.String(), "64") {
			name = "Int"
		}

		if name == "" {
			continue
		}

		goKind := protoutil.GoKindMap[protoKind]
		if goKind == 0 {
			continue
		}

		scalarTypeMap[protoKind] = &ScalarType{name: name, ProtoName: protoKind.String(), GoName: goKind.String()}
	}
}

type ScalarType struct {
	name      string
	ProtoName string
	GoName    string
}

func newScalarType(name string) *ScalarType { return &ScalarType{name: name} }
func (t *ScalarType) Name() string          { return t.name }
func (t *ScalarType) IsNullable() bool      { return false }
func (t *ScalarType) IsList() bool          { return false }
func (t *ScalarType) TypeAST() *ast.Type    { return ast.NonNullNamedType(t.Name(), nil) }

type WrappedScalarType struct {
	*ScalarType
	Proto *protogen.Message
}

func (t *WrappedScalarType) ProtoDescriptor() protoreflect.Descriptor { return t.Proto.Desc }
func (t *WrappedScalarType) GoIdent() protogen.GoIdent                { return t.Proto.GoIdent }

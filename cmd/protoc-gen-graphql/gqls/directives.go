package gqls

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
)

func messageDirectivesAST(m *protogen.Message) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "protobuf", Arguments: ast.ArgumentList{
			{Name: "type", Value: &ast.Value{Raw: string(m.Desc.FullName()), Kind: ast.StringValue}},
		}},
		{Name: "goModel", Arguments: ast.ArgumentList{
			{Name: "model", Value: &ast.Value{
				Raw:  fmt.Sprintf("%s.%s", string(m.GoIdent.GoImportPath), m.GoIdent.GoName),
				Kind: ast.StringValue},
			},
		}},
	}
}

func enumDirectivesAST(e *protogen.Enum) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "protobuf", Arguments: ast.ArgumentList{
			{Name: "type", Value: &ast.Value{Raw: string(e.Desc.FullName()), Kind: ast.StringValue}},
		}},
		{Name: "goModel", Arguments: ast.ArgumentList{
			{Name: "model", Value: &ast.Value{
				Raw:  fmt.Sprintf("%s.%s", string(e.GoIdent.GoImportPath), e.GoIdent.GoName),
				Kind: ast.StringValue},
			},
		}},
	}
}

func fieldDirectivesAST(f *protogen.Field) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "goField", Arguments: ast.ArgumentList{
			{Name: "name", Value: &ast.Value{Raw: f.GoName, Kind: ast.StringValue}},
			{Name: "forceResolver", Value: &ast.Value{Raw: "false", Kind: ast.BooleanValue}},
		}},
	}
}

func inputFieldDirectivesAST(f *protogen.Field) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "goField", Arguments: ast.ArgumentList{
			{Name: "name", Value: &ast.Value{Raw: f.GoName, Kind: ast.StringValue}},
		}},
	}
}

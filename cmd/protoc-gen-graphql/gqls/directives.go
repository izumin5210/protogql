package gqls

import (
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func rpcDirectivesAST(md protoreflect.MethodDescriptor) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "grpc", Arguments: ast.ArgumentList{
			{Name: "service", Value: &ast.Value{Raw: string(md.FullName().Parent()), Kind: ast.StringValue}},
			{Name: "rpc", Value: &ast.Value{Raw: string(md.Name()), Kind: ast.StringValue}},
		}},
	}
}

func definitionDelectivesAST(d protoreflect.Descriptor) ast.DirectiveList {
	return ast.DirectiveList{
		{Name: "protobuf", Arguments: ast.ArgumentList{
			{Name: "type", Value: &ast.Value{Raw: string(d.FullName()), Kind: ast.StringValue}},
		}},
	}
}

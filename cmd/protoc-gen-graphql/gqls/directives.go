package gqls

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func messageDirectivesAST(m *protogen.Message) ast.DirectiveList {
	return ast.DirectiveList{
		protobufTypeDirectiveAST(m.Desc, m.GoIdent),
	}
}

func enumDirectivesAST(e *protogen.Enum) ast.DirectiveList {
	return ast.DirectiveList{
		protobufTypeDirectiveAST(e.Desc, e.GoIdent),
	}
}

func fieldDirectivesAST(f *protogen.Field, typ Type) ast.DirectiveList {
	return ast.DirectiveList{
		protobufFieldDirectiveAST(f, typ),
	}
}

func inputFieldDirectivesAST(f *protogen.Field, typ Type) ast.DirectiveList {
	return ast.DirectiveList{
		protobufFieldDirectiveAST(f, typ),
	}
}

func protobufTypeDirectiveAST(desc protoreflect.Descriptor, goIdent protogen.GoIdent) *ast.Directive {
	return &ast.Directive{
		Name: "proto",
		Arguments: ast.ArgumentList{
			{Name: "fullName", Value: &ast.Value{Raw: string(desc.FullName()), Kind: ast.StringValue}},
			{Name: "package", Value: &ast.Value{Raw: string(desc.ParentFile().Package()), Kind: ast.StringValue}},
			{Name: "name", Value: &ast.Value{Raw: string(desc.Name()), Kind: ast.StringValue}},
			{Name: "goPackage", Value: &ast.Value{Raw: string(goIdent.GoImportPath), Kind: ast.StringValue}},
			{Name: "goName", Value: &ast.Value{Raw: goIdent.GoName, Kind: ast.StringValue}},
		},
	}
}

func protobufFieldDirectiveAST(f *protogen.Field, typ Type) *ast.Directive {
	var protoType, goType string
	switch typ := UnwrapType(typ).(type) {
	case *ScalarType:
		protoType = typ.ProtoName
		goType = typ.GoName
	case *ObjectType:
		protoType = string(typ.Proto.Desc.FullName())
		goType = fmt.Sprintf("%s.%s", typ.Proto.GoIdent.GoImportPath, typ.Proto.GoIdent.GoName)
	case *InputObjectType:
		protoType = string(typ.base.Proto.Desc.FullName())
		goType = fmt.Sprintf("%s.%s", typ.base.Proto.GoIdent.GoImportPath, typ.base.Proto.GoIdent.GoName)
	case *EnumType:
		protoType = string(typ.Proto.Desc.FullName())
		goType = fmt.Sprintf("%s.%s", typ.Proto.GoIdent.GoImportPath, typ.Proto.GoIdent.GoName)
	default:
		panic("unreachable")
	}
	return &ast.Directive{
		Name: "protoField",
		Arguments: ast.ArgumentList{
			{Name: "name", Value: &ast.Value{Raw: string(f.Desc.Name()), Kind: ast.StringValue}},
			{Name: "type", Value: &ast.Value{Raw: protoType, Kind: ast.StringValue}},
			{Name: "goName", Value: &ast.Value{Raw: f.GoName, Kind: ast.StringValue}},
			{Name: "goType", Value: &ast.Value{Raw: goType, Kind: ast.StringValue}},
		},
	}
}

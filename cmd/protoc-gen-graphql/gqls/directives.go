package gqls

import (
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
	if f.Oneof != nil {
		return ast.DirectiveList{
			protobufFieldDirectiveAST(f.Oneof.GoName, f.Oneof.Desc, typ),
		}
	}
	return ast.DirectiveList{
		protobufFieldDirectiveAST(f.GoName, f.Desc, typ),
	}
}

func inputFieldDirectivesAST(f *protogen.Field, typ Type) ast.DirectiveList {
	return ast.DirectiveList{
		protobufFieldDirectiveAST(f.GoName, f.Desc, typ),
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

func protobufFieldDirectiveAST(goName string, desc protoreflect.Descriptor, typ Type) *ast.Directive {
	var protoType, goTypeName, goTypePackage string

	switch typ := UnwrapType(typ).(type) {
	case *ScalarType:
		protoType = typ.ProtoName
		goTypeName = typ.GoName
	case ProtoType:
		protoType = string(typ.ProtoDescriptor().FullName())
		goTypeName = typ.GoIdent().GoName
		goTypePackage = string(typ.GoIdent().GoImportPath)
	default:
		panic("unreachable")
	}

	d := &ast.Directive{
		Name: "protoField",
		Arguments: ast.ArgumentList{
			{Name: "name", Value: &ast.Value{Raw: string(desc.Name()), Kind: ast.StringValue}},
			{Name: "type", Value: &ast.Value{Raw: protoType, Kind: ast.StringValue}},
			{Name: "goName", Value: &ast.Value{Raw: goName, Kind: ast.StringValue}},
			{Name: "goTypeName", Value: &ast.Value{Raw: goTypeName, Kind: ast.StringValue}},
		},
	}
	if goTypePackage != "" {
		d.Arguments = append(d.Arguments,
			&ast.Argument{Name: "goTypePackage", Value: &ast.Value{Raw: goTypePackage, Kind: ast.StringValue}})
	}
	return d
}

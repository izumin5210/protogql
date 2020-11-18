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

func oneofDirectivesAST(o *protogen.Oneof, gqlTypeNames []string) ast.DirectiveList {
	d := protobufTypeDirectiveAST(o.Desc, o.GoIdent)

	oneofFields := make(ast.ChildValueList, len(o.Fields))
	for i, f := range o.Fields {
		oneofFields[i] = &ast.ChildValue{Value: &ast.Value{
			Children: ast.ChildValueList{
				{Name: "name", Value: &ast.Value{Raw: gqlTypeNames[i], Kind: ast.StringValue}},
				{Name: "goName", Value: &ast.Value{Raw: f.GoIdent.GoName, Kind: ast.StringValue}},
			},
			Kind: ast.ObjectValue,
		}}
	}

	d.Arguments = append(d.Arguments, &ast.Argument{
		Name: "oneof",
		Value: &ast.Value{
			Children: ast.ChildValueList{
				{Name: "fields", Value: &ast.Value{Children: oneofFields, Kind: ast.ListValue}},
			},
			Kind: ast.ObjectValue,
		},
	})

	return ast.DirectiveList{d}
}

func enumDirectivesAST(e *protogen.Enum) ast.DirectiveList {
	return ast.DirectiveList{
		protobufTypeDirectiveAST(e.Desc, e.GoIdent),
	}
}

func oneofFieldDirectivesAST(o *protogen.Oneof, typ Type) ast.DirectiveList {
	return ast.DirectiveList{
		&ast.Directive{
			Name: "protoField",
			Arguments: append(
				append(ast.ArgumentList{}, protobufFieldDirectiveNameArgs(o.GoName, o.Desc.Name(), o)...),
				append(ast.ArgumentList{}, protobufFieldDirectiveTypeArgs(typ)...)...,
			),
		},
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
	return &ast.Directive{
		Name: "protoField",
		Arguments: append(
			append(ast.ArgumentList{}, protobufFieldDirectiveNameArgs(f.GoName, f.Desc.Name(), f.Oneof)...),
			append(ast.ArgumentList{}, protobufFieldDirectiveTypeArgs(typ)...)...,
		),
	}
}

func protobufFieldDirectiveTypeArgs(typ Type) []*ast.Argument {
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

	args := []*ast.Argument{
		{Name: "type", Value: &ast.Value{Raw: protoType, Kind: ast.StringValue}},
		{Name: "goTypeName", Value: &ast.Value{Raw: goTypeName, Kind: ast.StringValue}},
	}
	if goTypePackage != "" {
		args = append(args,
			&ast.Argument{Name: "goTypePackage", Value: &ast.Value{Raw: goTypePackage, Kind: ast.StringValue}})
	}

	return args
}

func protobufFieldDirectiveNameArgs(goName string, name protoreflect.Name, o *protogen.Oneof) []*ast.Argument {
	args := []*ast.Argument{
		{Name: "name", Value: &ast.Value{Raw: string(name), Kind: ast.StringValue}},
		{Name: "goName", Value: &ast.Value{Raw: goName, Kind: ast.StringValue}},
	}
	if o != nil {
		args = append(args,
			&ast.Argument{Name: "oneofName", Value: &ast.Value{Raw: string(o.Desc.Name()), Kind: ast.StringValue}},
			&ast.Argument{Name: "oneofGoName", Value: &ast.Value{Raw: o.GoName, Kind: ast.StringValue}},
		)
	}
	return args
}

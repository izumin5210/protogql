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

func fieldDirectivesAST(f *protogen.Field, typ Type) ast.DirectiveList {
	if f.Oneof != nil {
		return ast.DirectiveList{
			protobufOneofFieldDirectiveAST(f, typ),
		}
	}

	return ast.DirectiveList{
		protobufFieldDirectiveAST(f.GoName, f.Desc, typ),
	}
}

func inputFieldDirectivesAST(f *protogen.Field, typ Type) ast.DirectiveList {
	if f.Oneof != nil {
		return ast.DirectiveList{
			protobufOneofFieldDirectiveAST(f, typ),
		}
	}

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

func protobufOneofFieldDirectiveAST(f *protogen.Field, typ Type) *ast.Directive {
	d := protobufFieldDirectiveAST(f.Oneof.GoName, f.Oneof.Desc, typ)
	d.Arguments = append(d.Arguments,
		&ast.Argument{Name: "oneofName", Value: &ast.Value{Raw: string(f.Oneof.Desc.Name()), Kind: ast.StringValue}},
		&ast.Argument{Name: "oneofGoName", Value: &ast.Value{Raw: f.Oneof.GoName, Kind: ast.StringValue}},
	)
	return d
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

package types

import (
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type ObjectHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (o *ObjectHasProto) GoWrapperTypeName() string {
	return o.GoTypeName() + "_Proto"
}

func (o *ObjectHasProto) GoTypeName() string {
	return o.def.Name
}

func (o *ObjectHasProto) Godoc() string {
	return goutil.ToComment(o.def.Description)
}

func (o *ObjectHasProto) Fields() ([]*FieldHasProto, error) {
	fields := make([]*FieldHasProto, len(o.def.Fields))

	for i, f := range o.def.Fields {
		fields[i] = &FieldHasProto{gql: f, object: o}
	}

	return fields, nil
}

func (o *ObjectHasProto) FuncNameFromProto() string {
	return o.GoTypeName() + "FromProto"
}

func (o *ObjectHasProto) FuncNameFromRepeatedProto() string {
	return o.GoTypeName() + "ListFromRepeatedProto"
}

func (o *ObjectHasProto) FuncNameToProto() string {
	return o.GoTypeName() + "ToProto"
}

func (o *ObjectHasProto) FuncNameToRepeatedProto() string {
	return o.GoTypeName() + "ListToRepeatedProto"
}

func (o *ObjectHasProto) CodegenObject() *codegen.Object {
	return o.registry.FindObjectOrInput(o.def)
}

type FieldHasProto struct {
	gql    *ast.FieldDefinition
	object *ObjectHasProto
}

func (f *FieldHasProto) GoFieldName() string {
	return templates.ToGo(f.gql.Name)
}

func (f *FieldHasProto) Godoc() string {
	return goutil.ToComment(f.gql.Description)
}

func (f *FieldHasProto) GoFieldTypeDefinition() string {
	var b strings.Builder

	if f.isList() {
		b.WriteString("[]")
	}

	if !f.isGoBuiltinType() {
		b.WriteString("*")
	}

	switch typ := f.object.registry.FindType(f.gql.Type.Name()).(type) {
	case ProtoType:
		b.WriteString(typ.PbGoTypeName())
	case *ObjectHasProto:
		b.WriteString(typ.GoWrapperTypeName())
	default:
		for _, field := range f.object.CodegenObject().Fields {
			if field.Name == f.gql.Name {
				b.Reset()
				b.WriteString(templates.CurrentImports.LookupType(field.TypeReference.GO))
				break
			}
		}
	}

	return b.String()
}

func (f *FieldHasProto) FromProtoStatement(receiver string) string {
	var b strings.Builder

	if typ := f.object.registry.FindProtoLikeType(f.gql.Type.Name()); typ != nil {
		if f.isList() {
			b.WriteString(typ.FuncNameFromRepeatedProto())
		} else {
			b.WriteString(typ.FuncNameFromProto())
		}
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	} else {
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
	}

	return b.String()
}

func (f *FieldHasProto) ToProtoStatement(receiver string) string {
	var b strings.Builder

	if typ := f.object.registry.FindProtoLikeType(f.gql.Type.Name()); typ != nil {
		if f.isList() {
			b.WriteString(typ.FuncNameToRepeatedProto())
		} else {
			b.WriteString(typ.FuncNameToProto())
		}
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	} else {
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
	}

	return b.String()
}

func (f *FieldHasProto) isList() bool {
	return gqlutil.IsListType(f.gql.Type)
}

func (f *FieldHasProto) isGoBuiltinType() bool {
	return gqlutil.IsBuiltinType(f.gql.Type)
}

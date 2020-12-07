package types

import (
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

var (
	_ ProtoType = (*ObjectHasProto)(nil)
)

type ObjectHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (u *ObjectHasProto) IsFromProto() bool { return false }

func (o *ObjectHasProto) GQLName() string {
	return o.def.Name
}

func (o *ObjectHasProto) ProtoGoType() GoType {
	return newGoModelType(o.GoType().Name() + "_Proto")
}

func (o *ObjectHasProto) GoType() GoType {
	return newGoModelType(o.def.Name)
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

func (o *ObjectHasProto) ImplementedInterfaces() ([]Type, error) {
	types := make([]Type, len(o.def.Interfaces))

	for i, ifName := range o.def.Interfaces {
		typ := o.registry.FindInterfaceType(ifName)
		if typ == nil {
			return nil, errors.Errorf("interface %s was not found", ifName)
		}
		types[i] = typ
	}

	return types, nil
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

	switch typ := f.object.registry.FindType(f.gql.Type.Name()).(type) {
	case ProtoType:
		b.WriteString(typ.ProtoGoType().TypeString())
	default:
		if !f.isGoBuiltinType() {
			b.WriteString("*")
		}
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

	if typ := f.object.registry.FindProtoType(f.gql.Type.Name()); typ != nil {
		if f.isList() {
			b.WriteString(FromRepeatedProtoFuncName(typ))
		} else {
			b.WriteString(FromProtoFuncName(typ))
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

	if typ := f.object.registry.FindProtoType(f.gql.Type.Name()); typ != nil {
		if f.isList() {
			b.WriteString(ToRepeatedProtoFuncName(typ))
		} else {
			b.WriteString(ToProtoFuncName(typ))
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

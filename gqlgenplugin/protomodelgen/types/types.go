package types

import (
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
)

type ProtoType interface {
	Type
	ProtoGoType() GoType
	IsFromProto() bool
}

type Type interface {
	GQLName() string
	GoType() GoType
}

type UnionMemberType interface {
	Type() Type
}

type GoType interface {
	Name() string
	Package() string
	Pointer() bool
	QualifiedIdent() string
	TypeString() string
}

type goType struct {
	pkg  string
	name string
	ptr  bool
}

func newGoModelType(name string) GoType          { return &goType{name: name, ptr: true} }
func newGoModelInterfaceType(name string) GoType { return &goType{name: name} }
func newGoType(pkg, name string) GoType          { return &goType{pkg: pkg, name: name, ptr: true} }
func newGoInterfaceType(pkg, name string) GoType { return &goType{pkg: pkg, name: name} }

func (t *goType) Name() string           { return t.name }
func (t *goType) Package() string        { return t.pkg }
func (t *goType) Pointer() bool          { return t.ptr }
func (t *goType) TypeString() string     { return TypeString(t) }
func (t *goType) QualifiedIdent() string { return QualifiedIdent(t) }

type wrappedGoType struct {
	GoType
	name string
}

func wrapGoType(typ GoType, name string) GoType {
	return &wrappedGoType{GoType: typ, name: name}
}

func (t *wrappedGoType) Name() string           { return t.name }
func (t *wrappedGoType) TypeString() string     { return TypeString(t) }
func (t *wrappedGoType) QualifiedIdent() string { return QualifiedIdent(t) }

func QualifiedIdent(t GoType) string {
	var b strings.Builder
	if t.Package() != "" {
		b.WriteString(templates.CurrentImports.Lookup(t.Package()))
		b.WriteString(".")
	}
	b.WriteString(t.Name())
	return b.String()
}

func TypeString(t GoType) string {
	var b strings.Builder
	if t.Pointer() {
		b.WriteString("*")
	}
	b.WriteString(QualifiedIdent(t))
	return b.String()
}

func GoWrapperType(t Type) GoType {
	switch t := t.(type) {
	case ProtoType:
		return t.ProtoGoType()
	default:
		return t.GoType()
	}
}

func UnwrapStatement(t Type, varName string) string {
	switch t := t.(type) {
	case ProtoType:
		return fmt.Sprintf("%s(%s)", ToProtoFuncName(t), varName)
	default:
		return varName
	}
}

func IsFromProto(t Type) bool {
	pt, ok := t.(ProtoType)
	return ok && pt.IsFromProto()
}

func FuncCallStatement(funcname string, args ...string) string {
	return fmt.Sprintf("%s(%s)", funcname, strings.Join(args, ", "))
}

func FromProtoFuncName(typ Type) string {
	return fromProtoFuncName(typ, false)
}

func FromRepeatedProtoFuncName(typ Type) string {
	return fromProtoFuncName(typ, true)
}

func fromProtoFuncName(typ Type, repeated bool) string {
	if repeated {
		return typ.GQLName() + "ListFromRepeatedProto"
	}
	return typ.GQLName() + "FromProto"
}

func FromProtoFuncSignature(typ ProtoType) string {
	return fromProtoFuncSignature(typ, false)
}

func FromRepeatedProtoFuncSignature(typ ProtoType) string {
	return fromProtoFuncSignature(typ, true)
}

func fromProtoFuncSignature(typ ProtoType, repeated bool) string {
	return mapModelFuncSignature(fromProtoFuncName(typ, repeated), typ.ProtoGoType(), typ.GoType(), repeated)
}

func ToProtoFuncName(typ Type) string {
	return toProtoFuncName(typ, false)
}

func ToRepeatedProtoFuncName(typ Type) string {
	return toProtoFuncName(typ, true)
}

func toProtoFuncName(typ Type, repeated bool) string {
	if repeated {
		return typ.GQLName() + "ListToRepeatedProto"
	}
	return typ.GQLName() + "ToProto"
}

func ToProtoFuncSignature(typ ProtoType) string {
	return toProtoFuncSignature(typ, false)
}

func ToRepeatedProtoFuncSignature(typ ProtoType) string {
	return toProtoFuncSignature(typ, true)
}

func toProtoFuncSignature(typ ProtoType, repeated bool) string {
	return mapModelFuncSignature(toProtoFuncName(typ, repeated), typ.GoType(), typ.ProtoGoType(), repeated)
}

func mapModelFuncSignature(name string, reqType, respType GoType, repeated bool) string {
	var mod string
	if repeated {
		mod = "[]"
	}
	return fmt.Sprintf("%s(in %s%s) %s%s", name, mod, reqType.TypeString(), mod, respType.TypeString())
}

package protomodelgen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type Plugin struct {
}

func New() *Plugin {
	return new(Plugin)
}

var (
	_ plugin.ConfigMutator = (*Plugin)(nil)
	_ plugin.CodeGenerator = (*Plugin)(nil)

	funcs = template.FuncMap{
		"goWrapperTypeName": func(t Type) string {
			switch t := t.(type) {
			case ProtoType:
				return t.PbGoTypeName()
			case ProtoWrapperType:
				return t.GoWrapperTypeName()
			default:
				return t.GoTypeName()
			}
		},
		"unwrapStatement": func(t Type, varName string) string {
			switch t := t.(type) {
			case ProtoLikeType:
				return fmt.Sprintf("%s(%s)", t.FuncNameToProto(), varName)
			default:
				return varName
			}
		},
		"isProtoType": func(t Type) bool {
			_, ok := t.(ProtoType)
			return ok
		},
	}
)

func (p *Plugin) Name() string { return "protomodelgen" }

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	reg, err := CreateRegistryFromSchema(cfg.Schema)
	if err != nil {
		return err
	}

	models := []ProtoType{}
	for _, t := range reg.ObjectsFromProto() {
		models = append(models, t)
	}
	for _, t := range reg.EnumsFromProto() {
		models = append(models, t)
	}
	for _, t := range reg.UnionsFromProto() {
		models = append(models, t)
	}

	for _, typ := range models {
		cfg.Models.Add(typ.GoTypeName(), cfg.Model.ImportPath()+"."+typ.GoTypeName())
	}
	for _, obj := range reg.ObjectsFromProto() {
		fields, err := obj.Fields()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, f := range fields {
			if f.proto != nil {
				continue
			}
			if cfg.Models[obj.def.Name].Fields == nil {
				cfg.Models[obj.def.Name] = config.TypeMapEntry{
					Model:  cfg.Models[obj.def.Name].Model,
					Fields: map[string]config.TypeMapField{},
				}
			}
			cfg.Models[obj.def.Name].Fields[f.gql.Name] = config.TypeMapField{FieldName: f.gql.Name, Resolver: true}
		}
	}

	for _, enum := range reg.EnumsFromProto() {
		cfg.Models.Add(enum.def.Name, cfg.Model.ImportPath()+"."+enum.GoTypeName())
	}

	cfg.Directives["grpc"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["proto"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["protoField"] = config.DirectiveConfig{SkipRuntime: true}

	for _, name := range []string{"Int", "ID"} {
		model := cfg.Models[name]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.Uint32")
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.Uint64")
		cfg.Models[name] = model
	}

	{
		model := cfg.Models["Int"]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.Int32Value")
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.Int64Value")
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.UInt32Value")
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.UInt64Value")
		cfg.Models["Int"] = model
	}
	{
		model := cfg.Models["Float"]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.FloatValue")
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.DoubleValue")
		cfg.Models["Float"] = model
	}
	{
		model := cfg.Models["Boolean"]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.BoolValue")
		cfg.Models["Boolean"] = model
	}
	{
		model := cfg.Models["String"]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.StringValue")
		cfg.Models["String"] = model
	}
	{
		model := cfg.Models["DateTime"]
		model.Model = append(model.Model, "github.com/izumin5210/protogql/gqlruntime/types.Timestamp")
		cfg.Models["DateTime"] = model
	}

	return templates.Render(templates.Options{
		PackageName:     cfg.Model.Package,
		Filename:        filepath.Join(cfg.Model.Dir(), "protomodels_gen.go"),
		Data:            reg,
		Funcs:           funcs,
		GeneratedHeader: true,
		Packages:        cfg.Packages,
	})
}

func (p *Plugin) GenerateCode(data *codegen.Data) error {
	reg, err := CreateRegistry(data)
	if err != nil {
		return err
	}

	return templates.Render(templates.Options{
		PackageName:     data.Config.Model.Package,
		Filename:        filepath.Join(data.Config.Model.Dir(), "protomodels_gen.go"),
		Data:            reg,
		GeneratedHeader: true,
		Funcs:           funcs,
		Packages:        data.Config.Packages,
	})
}

type Registry struct {
	objectsFromProto map[string]*ObjectFromProto
	objectsHasProto  map[string]*ObjectHasProto
	plainObjects     map[string]*PlainObject
	enumsFromProto   map[string]*EnumFromProto
	unionsFromProto  map[string]*UnionFromProto
	unionsHasProto   map[string]*UnionHasProto
	data             *codegen.Data
}

func CreateRegistry(data *codegen.Data) (*Registry, error) {
	return createRegistry(data, data.Schema)
}

func CreateRegistryFromSchema(schema *ast.Schema) (*Registry, error) {
	return createRegistry(nil, schema)
}

func createRegistry(data *codegen.Data, schema *ast.Schema) (*Registry, error) {
	reg := &Registry{
		objectsFromProto: map[string]*ObjectFromProto{},
		objectsHasProto:  map[string]*ObjectHasProto{},
		plainObjects:     map[string]*PlainObject{},
		enumsFromProto:   map[string]*EnumFromProto{},
		unionsFromProto:  map[string]*UnionFromProto{},
		unionsHasProto:   map[string]*UnionHasProto{},
		data:             data,
	}

	for _, def := range schema.Types {
		if strings.HasPrefix(def.Name, "__") {
			continue
		}
		if q, m := schema.Query, schema.Mutation; (q != nil && def.Name == q.Name) || (m != nil && def.Name == m.Name) {
			continue
		}

		switch def.Kind {
		case ast.Object, ast.InputObject:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.objectsFromProto[def.Name] = &ObjectFromProto{def: def, proto: proto, registry: reg}
			} else if ok, err := gqlutil.HasProto(def, schema.Types); err == nil && ok {
				reg.objectsHasProto[def.Name] = &ObjectHasProto{def: def, registry: reg}
			} else {
				reg.plainObjects[def.Name] = &PlainObject{def: def}
			}

		case ast.Enum:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.enumsFromProto[def.Name] = &EnumFromProto{def: def, proto: proto}
			} else {
				panic("Plain GraphQL Enums is not supported yet")
			}

		case ast.Scalar:
			// no-op

		case ast.Union:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.unionsFromProto[def.Name] = &UnionFromProto{def: def, proto: proto, registry: reg}
			} else if ok, err := gqlutil.HasProto(def, schema.Types); err == nil && ok {
				reg.unionsHasProto[def.Name] = &UnionHasProto{def: def, registry: reg}
			} else {
				panic("Plain GraphQL Unions is not supported yet")
			}

		default:
			// TODO: not implemented
			panic(fmt.Errorf("%s is not supported yet", def.Kind))
		}
	}

	return reg, nil
}

func (r *Registry) FindType(name string) Type {
	if typ := r.FindProtoLikeType(name); typ != nil {
		return typ
	}
	if obj, ok := r.plainObjects[name]; ok {
		return obj
	}

	return nil
}

func (r *Registry) FindProtoType(name string) ProtoType {
	if obj, ok := r.objectsFromProto[name]; ok {
		return obj
	}
	if enum, ok := r.enumsFromProto[name]; ok {
		return enum
	}
	if union, ok := r.unionsFromProto[name]; ok {
		return union
	}

	return nil
}

func (r *Registry) FindProtoLikeType(name string) ProtoLikeType {
	if typ := r.FindProtoType(name); typ != nil {
		return typ
	}
	if obj, ok := r.objectsHasProto[name]; ok {
		return obj
	}
	if union, ok := r.unionsHasProto[name]; ok {
		return union
	}

	return nil
}

func (r *Registry) FindObjectOrInput(def *ast.Definition) *codegen.Object {
	if def.Kind == ast.InputObject {
		return r.data.Inputs.ByName(def.Name)
	}
	return r.data.Objects.ByName(def.Name)
}

func (r *Registry) ObjectsFromProto() []*ObjectFromProto {
	objs := make([]*ObjectFromProto, 0, len(r.objectsFromProto))
	for _, o := range r.objectsFromProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GoTypeName() < objs[j].GoTypeName() })

	return objs
}

func (r *Registry) ObjectsHasProto() []*ObjectHasProto {
	// FIXME
	if r.data == nil {
		return []*ObjectHasProto{}
	}

	objs := make([]*ObjectHasProto, 0, len(r.objectsHasProto))
	for _, o := range r.objectsHasProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GoTypeName() < objs[j].GoTypeName() })

	return objs
}

func (r *Registry) EnumsFromProto() []*EnumFromProto {
	enums := make([]*EnumFromProto, 0, len(r.enumsFromProto))
	for _, e := range r.enumsFromProto {
		enums = append(enums, e)
	}

	sort.Slice(enums, func(i, j int) bool { return enums[i].GoTypeName() < enums[j].GoTypeName() })

	return enums
}

func (r *Registry) UnionsFromProto() []*UnionFromProto {
	unions := make([]*UnionFromProto, 0, len(r.unionsFromProto))
	for _, u := range r.unionsFromProto {
		unions = append(unions, u)
	}

	sort.Slice(unions, func(i, j int) bool { return unions[i].GoTypeName() < unions[j].GoTypeName() })

	return unions
}

func (r *Registry) UnionsHasProto() []*UnionHasProto {
	unions := make([]*UnionHasProto, 0, len(r.unionsHasProto))
	for _, u := range r.unionsHasProto {
		unions = append(unions, u)
	}

	sort.Slice(unions, func(i, j int) bool { return unions[i].GoTypeName() < unions[j].GoTypeName() })

	return unions
}

type ProtoType interface {
	ProtoLikeType
	PbGoTypeName() string
}

type ProtoWrapperType interface {
	ProtoLikeType
	GoWrapperTypeName() string
}

type ProtoLikeType interface {
	Type
	FuncNameFromProto() string
	FuncNameFromRepeatedProto() string
	FuncNameToProto() string
	FuncNameToRepeatedProto() string
}

type Type interface {
	GoTypeName() string
}

type PlainObject struct {
	def *ast.Definition
}

func (o *PlainObject) GoTypeName() string {
	return o.def.Name
}

type ObjectFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (o *ObjectFromProto) GoTypeName() string {
	return o.def.Name
}

func (o *ObjectFromProto) Godoc() string {
	return goutil.ToComment(o.def.Description)
}

func (o *ObjectFromProto) Fields() ([]*FieldFromProto, error) {
	fields := make([]*FieldFromProto, len(o.def.Fields))

	for i, f := range o.def.Fields {
		proto, err := gqlutil.ExtractProtoFieldDirective(f.Directives)
		if err != nil {
			return nil, errors.Wrapf(err, "%s has invalid directive", f.Name)
		}
		fields[i] = &FieldFromProto{gql: f, proto: proto, object: o}
	}

	return fields, nil
}

func (o *ObjectFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(o.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(o.proto.GoName)

	return b.String()
}

func (o *ObjectFromProto) FuncNameFromProto() string {
	return o.GoTypeName() + "FromProto"
}

func (o *ObjectFromProto) FuncNameFromRepeatedProto() string {
	return o.GoTypeName() + "ListFromRepeatedProto"
}

func (o *ObjectFromProto) FuncNameToProto() string {
	return o.GoTypeName() + "ToProto"
}

func (o *ObjectFromProto) FuncNameToRepeatedProto() string {
	return o.GoTypeName() + "ListToRepeatedProto"
}

type FieldFromProto struct {
	gql    *ast.FieldDefinition
	proto  *gqlutil.ProtoFieldDirective
	object *ObjectFromProto
}

func (f *FieldFromProto) GoFieldName() string {
	return templates.ToGo(f.gql.Name)
}

func (f *FieldFromProto) PbGoFieldName() string {
	return f.proto.GoName
}

func (f *FieldFromProto) Godoc() string {
	return goutil.ToComment(f.gql.Description)
}

func (f *FieldFromProto) GoFieldTypeDefinition() string {
	var b strings.Builder

	if f.isList() {
		b.WriteString("[]")
	}

	switch {
	case f.isGoBuiltinType():
		b.WriteString(f.proto.GoTypeName)
	case f.isProtoWellKnownType():
		b.WriteString("*")
		b.WriteString(templates.CurrentImports.Lookup(f.proto.GoTypePackage))
		b.WriteString(".")
		b.WriteString(f.proto.GoTypeName)
	default:
		typ := f.object.registry.FindType(f.gql.Type.Name())
		if _, ok := typ.(*UnionFromProto); !ok {
			b.WriteString("*")
		}
		b.WriteString(typ.GoTypeName())
	}

	return b.String()
}

func (f *FieldFromProto) FromProtoStatement(receiver string) string {
	if f.proto == nil {
		return ""
	}

	var b strings.Builder

	switch {
	case f.isGoBuiltinType(), f.isProtoWellKnownType():
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("()")
	case f.isList():
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameFromRepeatedProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("())")
	default:
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameFromProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".Get")
		b.WriteString(f.proto.GoName)
		b.WriteString("())")
	}

	return b.String()
}

func (f *FieldFromProto) ToProtoStatement(receiver string) string {
	if f.proto == nil {
		return ""
	}

	var b strings.Builder

	switch {
	case f.isGoBuiltinType(), f.isProtoWellKnownType():
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
	case f.isList():
		b.WriteString(f.ProtoType().FuncNameToRepeatedProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	default:
		b.WriteString(f.ProtoType().FuncNameToProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	}

	return b.String()
}

func (f *FieldFromProto) ProtoType() ProtoType {
	return f.object.registry.FindProtoType(f.gql.Type.Name())
}

func (f *FieldFromProto) IsOneof() bool {
	return f.proto != nil && f.proto.OneofName != ""
}

func (f *FieldFromProto) OneofMembers() []*UnionMemberFromProto {
	typ := f.object.registry.FindProtoType(f.gql.Type.Name())
	u, ok := typ.(*UnionFromProto)
	if !ok {
		return nil
	}

	return u.Members()
}

func (f *FieldFromProto) IsOneofMember() bool {
	return f.proto != nil && f.proto.OneofName != "" && f.proto.OneofName != f.proto.Name
}

func (f *FieldFromProto) isList() bool {
	return gqlutil.IsListType(f.gql.Type)
}

func (f *FieldFromProto) isGoBuiltinType() bool {
	if f.proto != nil {
		return f.proto.IsGoBuiltinType()
	}
	return gqlutil.IsBuiltinType(f.gql.Type)
}

func (f *FieldFromProto) isProtoWellKnownType() bool {
	return f.proto != nil && f.proto.IsWellKnownType()
}

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

type EnumFromProto struct {
	def   *ast.Definition
	proto *gqlutil.ProtoDirective
}

func (e *EnumFromProto) GoTypeName() string {
	return e.def.Name
}

func (e *EnumFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(e.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(e.proto.GoName)

	return b.String()
}

func (e *EnumFromProto) Godoc() string {
	return goutil.ToComment(e.def.Description)
}

func (e *EnumFromProto) FuncNameFromProto() string {
	return e.GoTypeName() + "FromProto"
}

func (e *EnumFromProto) FuncNameFromRepeatedProto() string {
	return e.GoTypeName() + "ListFromRepeatedProto"
}

func (e *EnumFromProto) FuncNameToProto() string {
	return e.GoTypeName() + "ToProto"
}

func (e *EnumFromProto) FuncNameToRepeatedProto() string {
	return e.GoTypeName() + "ListToRepeatedProto"
}

type UnionMemberType interface {
	Type() Type
	UnionType() Type
}

type UnionFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (u *UnionFromProto) GoTypeName() string {
	return u.def.Name
}

func (u *UnionFromProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionFromProto) FuncNameFromProto() string {
	return u.GoTypeName() + "FromProto"
}

func (u *UnionFromProto) PbGoTypeName() string              { panic("unreachable") }
func (u *UnionFromProto) FuncNameFromRepeatedProto() string { panic("unreachable") }
func (u *UnionFromProto) FuncNameToProto() string           { panic("unreachable") }
func (u *UnionFromProto) FuncNameToRepeatedProto() string   { panic("unreachable") }

func (u *UnionFromProto) Members() []*UnionMemberFromProto {
	members := make([]*UnionMemberFromProto, len(u.proto.Oneof.Fields))
	for i, f := range u.proto.Oneof.Fields {
		members[i] = &UnionMemberFromProto{union: u, proto: f, typ: u.registry.FindProtoType(f.Name)}
	}
	return members
}

type UnionMemberFromProto struct {
	union *UnionFromProto
	proto *gqlutil.ProtoDirectiveOneofField
	typ   ProtoType
}

func (m *UnionMemberFromProto) GoTypeName() string {
	return m.typ.GoTypeName()
}

func (m *UnionMemberFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(m.union.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(m.proto.GoName)

	return b.String()
}

func (m *UnionMemberFromProto) PbGoTypeFieldName() string {
	return m.proto.Name
}

func (m *UnionMemberFromProto) FuncNameFromProto() string {
	return m.typ.FuncNameFromProto()
}

func (m *UnionMemberFromProto) FuncNameToProto() string {
	return m.proto.GoName + "ToProto"
}

func (m *UnionMemberFromProto) Type() ProtoType {
	return m.typ
}

type UnionHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (u *UnionHasProto) GoTypeName() string {
	return u.def.Name
}

func (u *UnionHasProto) GoWrapperTypeName() string {
	return u.GoTypeName() + "_Proto"
}

func (u *UnionHasProto) FuncNameFromProto() string {
	return u.GoTypeName() + "FromProto"
}

func (u *UnionHasProto) FuncNameFromRepeatedProto() string {
	return u.GoTypeName() + "ListFromRepeatedProto"
}

func (u *UnionHasProto) FuncNameToProto() string {
	return u.GoTypeName() + "ToProto"
}

func (u *UnionHasProto) FuncNameToRepeatedProto() string {
	return u.GoTypeName() + "ListToRepeatedProto"
}

func (u *UnionHasProto) Godoc() string {
	return goutil.ToComment(u.def.Description)
}

func (u *UnionHasProto) Members() []UnionMemberType {
	members := make([]UnionMemberType, len(u.def.Types))
	for i, t := range u.def.Types {
		switch typ := u.registry.FindType(t).(type) {
		case ProtoLikeType:
			members[i] = &UnionMemberHasProto{typ: typ}
		case Type:
			members[i] = &PlainUnionMember{typ: typ}
		default:
			panic("unreachable")
		}
	}
	return members
}

func (u *UnionHasProto) ProtoMembers() []*UnionMemberHasProto {
	members := make([]*UnionMemberHasProto, 0, len(u.def.Types))
	for _, m := range u.Members() {
		if mp, ok := m.(*UnionMemberHasProto); ok {
			members = append(members, mp)
		}
	}
	return members
}

type PlainUnionMember struct {
	typ   Type
	union Type
}

func (m *PlainUnionMember) Type() Type {
	return m.typ
}

func (m *PlainUnionMember) UnionType() Type {
	return m.union
}

type UnionMemberHasProto struct {
	typ   ProtoLikeType
	union *UnionHasProto
}

func (m *UnionMemberHasProto) Type() Type {
	return m.typ
}

func (m *UnionMemberHasProto) UnionType() Type {
	return m.union
}

package protomodelgen

import (
	"fmt"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/codegen/gqlutil"
)

type Plugin struct {
}

func New() *Plugin {
	return new(Plugin)
}

var (
	_ plugin.ConfigMutator = (*Plugin)(nil)
	_ plugin.CodeGenerator = (*Plugin)(nil)
)

func (p *Plugin) Name() string { return "protomodelgen" }

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	reg, err := CreateRegistry(cfg.Schema)
	if err != nil {
		return err
	}

	for _, obj := range reg.ObjectsFromProto() {
		cfg.Models.Add(obj.def.Name, cfg.Model.ImportPath()+"."+obj.GoTypeName())
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
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.Uint32")
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.Uint64")
		cfg.Models[name] = model
	}

	{
		model := cfg.Models["Int"]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.Int32Value")
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.Int64Value")
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.UInt32Value")
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.UInt64Value")
		cfg.Models["Int"] = model
	}
	{
		model := cfg.Models["Float"]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.FloatValue")
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.DoubleValue")
		cfg.Models["Float"] = model
	}
	{
		model := cfg.Models["Boolean"]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.BoolValue")
		cfg.Models["Boolean"] = model
	}
	{
		model := cfg.Models["String"]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.StringValue")
		cfg.Models["String"] = model
	}
	{
		model := cfg.Models["DateTime"]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime/types.Timestamp")
		cfg.Models["DateTime"] = model
	}

	return templates.Render(templates.Options{
		PackageName:     cfg.Model.Package,
		Filename:        filepath.Join(cfg.Model.Dir(), "protomodels_gen.go"),
		Data:            reg,
		GeneratedHeader: true,
		Packages:        cfg.Packages,
	})
}

func (p *Plugin) GenerateCode(data *codegen.Data) error {
	reg, err := CreateRegistry(data.Schema)
	if err != nil {
		return err
	}

	reg.data = data

	return templates.Render(templates.Options{
		PackageName:     data.Config.Model.Package,
		Filename:        filepath.Join(data.Config.Model.Dir(), "protomodels_gen.go"),
		Data:            reg,
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

type Registry struct {
	objectsFromProto map[string]*NewObjectFromProto
	objectsHasProto  map[string]*NewObjectHasProto
	plainObjects     map[string]*NewPlainObject
	enumsFromProto   map[string]*NewEnumFromProto
	data             *codegen.Data
}

func CreateRegistry(schema *ast.Schema) (*Registry, error) {
	reg := &Registry{
		objectsFromProto: map[string]*NewObjectFromProto{},
		objectsHasProto:  map[string]*NewObjectHasProto{},
		plainObjects:     map[string]*NewPlainObject{},
		enumsFromProto:   map[string]*NewEnumFromProto{},
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
				reg.objectsFromProto[def.Name] = &NewObjectFromProto{def: def, proto: proto, registry: reg}
			} else if ok, err := gqlutil.HasProto(def, schema.Types); err == nil && ok {
				reg.objectsHasProto[def.Name] = &NewObjectHasProto{def: def, registry: reg}
			} else {
				reg.plainObjects[def.Name] = &NewPlainObject{def: def}
			}

		case ast.Enum:
			proto, err := gqlutil.ExtractProtoDirective(def.Directives)
			if err != nil {
				return nil, errors.Wrapf(err, "%s has invalid directive", def.Name)
			}
			if proto != nil {
				reg.enumsFromProto[def.Name] = &NewEnumFromProto{def: def, proto: proto}
			} else {
				panic("Plain GraphQL Enums is not supported yet")
			}

		case ast.Scalar:
			// no-op

		default:
			// TODO: not implemented
			panic(fmt.Errorf("%s is not supported yet", def.Kind))
		}
	}

	return reg, nil
}

func (r *Registry) FindType(name string) NewType {
	if obj, ok := r.objectsFromProto[name]; ok {
		return obj
	}
	if obj, ok := r.objectsHasProto[name]; ok {
		return obj
	}
	if obj, ok := r.plainObjects[name]; ok {
		return obj
	}
	if enum, ok := r.enumsFromProto[name]; ok {
		return enum
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

	return nil
}

func (r *Registry) FindProtoLikeType(name string) ProtoLikeType {
	if obj, ok := r.objectsFromProto[name]; ok {
		return obj
	}
	if obj, ok := r.objectsHasProto[name]; ok {
		return obj
	}
	if enum, ok := r.enumsFromProto[name]; ok {
		return enum
	}

	return nil
}

func (r *Registry) ObjectsFromProto() []*NewObjectFromProto {
	objs := make([]*NewObjectFromProto, 0, len(r.objectsFromProto))
	for _, o := range r.objectsFromProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GoTypeName() < objs[j].GoTypeName() })

	return objs
}

func (r *Registry) ObjectsHasProto() []*NewObjectHasProto {
	// FIXME
	if r.data == nil {
		return []*NewObjectHasProto{}
	}

	objs := make([]*NewObjectHasProto, 0, len(r.objectsHasProto))
	for _, o := range r.objectsHasProto {
		objs = append(objs, o)
	}

	sort.Slice(objs, func(i, j int) bool { return objs[i].GoTypeName() < objs[j].GoTypeName() })

	return objs
}

func (r *Registry) EnumsFromProto() []*NewEnumFromProto {
	enums := make([]*NewEnumFromProto, 0, len(r.enumsFromProto))
	for _, e := range r.enumsFromProto {
		enums = append(enums, e)
	}

	sort.Slice(enums, func(i, j int) bool { return enums[i].GoTypeName() < enums[j].GoTypeName() })

	return enums
}

type ProtoType interface {
	NewType
	ProtoLikeType
	PbGoTypeName() string
}
type ProtoLikeType interface {
	NewType
	FuncNameFromProto() string
	FuncNameFromRepeatedProto() string
	FuncNameToProto() string
	FuncNameToRepeatedProto() string
}

type NewType interface {
	GoTypeName() string
}

type NewPlainObject struct {
	def *ast.Definition
}

func (o *NewPlainObject) GoTypeName() string {
	return o.def.Name
}

type NewObjectFromProto struct {
	def      *ast.Definition
	proto    *gqlutil.ProtoDirective
	registry *Registry
}

func (o *NewObjectFromProto) GoTypeName() string {
	return o.def.Name
}

func (o *NewObjectFromProto) Fields() ([]*NewFieldFromProto, error) {
	fields := make([]*NewFieldFromProto, len(o.def.Fields))

	for i, f := range o.def.Fields {
		proto, err := gqlutil.ExtractProtoFieldDirective(f.Directives)
		if err != nil {
			return nil, errors.Wrapf(err, "%s has invalid directive", f.Name)
		}
		fields[i] = &NewFieldFromProto{gql: f, proto: proto, object: o}
	}

	return fields, nil
}

func (o *NewObjectFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(o.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(o.proto.GoName)

	return b.String()
}

func (o *NewObjectFromProto) FuncNameFromProto() string {
	return o.GoTypeName() + "FromProto"
}

func (o *NewObjectFromProto) FuncNameFromRepeatedProto() string {
	return o.GoTypeName() + "ListFromRepeatedProto"
}

func (o *NewObjectFromProto) FuncNameToProto() string {
	return o.GoTypeName() + "ToProto"
}

func (o *NewObjectFromProto) FuncNameToRepeatedProto() string {
	return o.GoTypeName() + "ListToRepeatedProto"
}

type NewFieldFromProto struct {
	gql    *ast.FieldDefinition
	proto  *gqlutil.ProtoFieldDirective
	object *NewObjectFromProto
}

func (f *NewFieldFromProto) GoFieldName() string {
	return templates.ToGo(f.gql.Name)
}

func (f *NewFieldFromProto) PbGoFieldName() string {
	return f.proto.GoName
}

func (f *NewFieldFromProto) GoFieldTypeDefinition() string {
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
		b.WriteString("*")
		typ := f.object.registry.FindType(f.gql.Type.Name())
		b.WriteString(typ.GoTypeName())
	}

	return b.String()
}

func (f *NewFieldFromProto) FromProtoStatement(receiver string) string {
	if f.proto == nil {
		return ""
	}

	var b strings.Builder

	switch {
	case f.isGoBuiltinType(), f.isProtoWellKnownType():
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.proto.GoName)
	case f.isList():
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameFromRepeatedProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.proto.GoName)
		b.WriteString(")")
	default:
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameFromProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.proto.GoName)
		b.WriteString(")")
	}

	return b.String()
}

func (f *NewFieldFromProto) ToProtoStatement(receiver string) string {
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
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameToRepeatedProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	default:
		typ := f.object.registry.FindProtoType(f.gql.Type.Name())
		b.WriteString(typ.FuncNameToProto())
		b.WriteString("(")
		b.WriteString(receiver)
		b.WriteString(".")
		b.WriteString(f.GoFieldName())
		b.WriteString(")")
	}

	return b.String()
}

func (f *NewFieldFromProto) isList() bool {
	return f.gql.Type.NamedType == ""
}

func (f *NewFieldFromProto) isGoBuiltinType() bool {
	if f.proto == nil {
		switch f.gql.Type.Name() {
		case "ID", "Int", "Float", "String", "Boolean":
			return true
		default:
			return false
		}
	}
	return strings.ToLower(f.proto.Type) == f.proto.Type
}

func (f *NewFieldFromProto) isProtoWellKnownType() bool {
	if f.proto == nil {
		return false
	}
	switch f.proto.Type {
	case "google.protobuf.Int32Value", "google.protobuf.Int64Value",
		"google.protobuf.UInt32Value", "google.protobuf.UInt64Value",
		"google.protobuf.FloatValue", "google.protobuf.DoubleValue",
		"google.protobuf.BoolValue",
		"google.protobuf.StringValue",
		"google.protobuf.Timestamp":
		return true
	}
	return false
}

type NewObjectHasProto struct {
	def      *ast.Definition
	registry *Registry
}

func (o *NewObjectHasProto) GoWrapperTypeName() string {
	return o.GoTypeName() + "_Proto"
}

func (o *NewObjectHasProto) GoTypeName() string {
	return o.def.Name
}

func (o *NewObjectHasProto) Fields() ([]*NewFieldHasProto, error) {
	fields := make([]*NewFieldHasProto, len(o.def.Fields))

	for i, f := range o.def.Fields {
		fields[i] = &NewFieldHasProto{gql: f, object: o}
	}

	return fields, nil
}

func (o *NewObjectHasProto) FuncNameFromProto() string {
	return o.GoTypeName() + "FromProto"
}

func (o *NewObjectHasProto) FuncNameFromRepeatedProto() string {
	return o.GoTypeName() + "ListFromRepeatedProto"
}

func (o *NewObjectHasProto) FuncNameToProto() string {
	return o.GoTypeName() + "ToProto"
}

func (o *NewObjectHasProto) FuncNameToRepeatedProto() string {
	return o.GoTypeName() + "ListToRepeatedProto"
}

type NewFieldHasProto struct {
	gql    *ast.FieldDefinition
	object *NewObjectHasProto
}

func (f *NewFieldHasProto) GoFieldName() string {
	return templates.ToGo(f.gql.Name)
}

func (f *NewFieldHasProto) GoFieldTypeDefinition() string {
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
	case *NewObjectHasProto:
		b.WriteString(typ.GoWrapperTypeName())
	default:
		byName := f.object.registry.data.Objects.ByName
		if f.object.def.Kind == ast.InputObject {
			byName = f.object.registry.data.Inputs.ByName
		}
		obj := byName(f.object.def.Name)
		for _, field := range obj.Fields {
			if field.Name == f.gql.Name {
				t := field.TypeReference.GO
				for {
					p, ok := t.(interface{ Elem() types.Type })
					if !ok {
						break
					}
					t = p.Elem()
				}
				b.WriteString(templates.CurrentImports.LookupType(t))
				break
			}
		}
	}

	return b.String()
}

func (f *NewFieldHasProto) FromProtoStatement(receiver string) string {
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

func (f *NewFieldHasProto) ToProtoStatement(receiver string) string {
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

func (f *NewFieldHasProto) isList() bool {
	return f.gql.Type.NamedType == ""
}

func (f *NewFieldHasProto) isGoBuiltinType() bool {
	switch f.gql.Type.Name() {
	case "ID", "Int", "Float", "String", "Boolean":
		return true
	default:
		return false
	}
}

type NewEnumFromProto struct {
	def   *ast.Definition
	proto *gqlutil.ProtoDirective
}

func (e *NewEnumFromProto) GoTypeName() string {
	return e.def.Name
}

func (e *NewEnumFromProto) PbGoTypeName() string {
	var b strings.Builder

	b.WriteString(templates.CurrentImports.Lookup(e.proto.GoPackage))
	b.WriteString(".")
	b.WriteString(e.proto.GoName)

	return b.String()
}

func (e *NewEnumFromProto) FuncNameFromProto() string {
	return e.GoTypeName() + "FromProto"
}

func (e *NewEnumFromProto) FuncNameFromRepeatedProto() string {
	return e.GoTypeName() + "ListFromRepeatedProto"
}

func (e *NewEnumFromProto) FuncNameToProto() string {
	return e.GoTypeName() + "ToProto"
}

func (e *NewEnumFromProto) FuncNameToRepeatedProto() string {
	return e.GoTypeName() + "ListToRepeatedProto"
}

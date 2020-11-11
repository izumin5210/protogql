package protomodelgen

import (
	"go/types"
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
	binding, err := createBinding(cfg.Schema)
	if err != nil {
		return err
	}

	for _, obj := range binding.ProtoObjects {
		cfg.Models.Add(obj.Name, cfg.Model.ImportPath()+"."+obj.Name)
		for _, f := range obj.Fields {
			if f.Proto != nil {
				continue
			}
			if cfg.Models[obj.Name].Fields == nil {
				cfg.Models[obj.Name] = config.TypeMapEntry{
					Model:  cfg.Models[obj.Name].Model,
					Fields: map[string]config.TypeMapField{},
				}
			}
			cfg.Models[obj.Name].Fields[f.Name] = config.TypeMapField{FieldName: f.Name, Resolver: true}
		}
	}
	for _, enum := range binding.Enums {
		cfg.Models.Add(enum.Name, cfg.Model.ImportPath()+"."+enum.Name)
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
		Data:            binding,
		GeneratedHeader: true,
		Packages:        cfg.Packages,
		Funcs: template.FuncMap{
			"findGQLFieldType": binding.FindGQLFieldType,
			"hasProto":         binding.HasProto,
		},
	})
}

func (p *Plugin) GenerateCode(data *codegen.Data) error {
	binding, err := createBinding(data.Schema)
	if err != nil {
		return errors.WithStack(err)
	}

	err = binding.PrepareObjectsHasProto()
	if err != nil {
		return errors.WithStack(err)
	}

	binding.data = data

	return templates.Render(templates.Options{
		PackageName:     data.Config.Model.Package,
		Filename:        filepath.Join(data.Config.Model.Dir(), "protomodels_gen.go"),
		Data:            binding,
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
		Funcs: template.FuncMap{
			"findGQLFieldType": binding.FindGQLFieldType,
			"hasProto":         binding.HasProto,
		},
	})
}

func createBinding(s *ast.Schema) (*Binding, error) {
	binding := new(Binding)
	binding.schema = s

	for _, typ := range s.Types {
		proto, err := gqlutil.ExtractProtoDirective(typ.Directives)
		if err != nil {
			return nil, errors.Wrapf(err, "%s has invalid directive", typ.Name)
		}
		if proto == nil {
			continue
		}

		switch typ.Kind {
		case ast.Object, ast.InputObject:
			obj, err := binding.newObject(typ)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			binding.ProtoObjects = append(binding.ProtoObjects, obj)

		case ast.Enum:
			enum := &Enum{Name: typ.Name, Proto: proto}
			for _, ev := range typ.EnumValues {
				enum.Values = append(enum.Values, &EnumValue{Name: ev.Name})
			}
			binding.Enums = append(binding.Enums, enum)
		}
	}

	sort.Slice(binding.ProtoObjects, func(i, j int) bool { return binding.ProtoObjects[i].Name < binding.ProtoObjects[j].Name })
	sort.Slice(binding.Enums, func(i, j int) bool { return binding.Enums[i].Name < binding.Enums[j].Name })

	return binding, nil
}

func (b *Binding) PrepareObjectsHasProto() error {
	for _, def := range b.schema.Types {
		if !(def.Kind == ast.Object || def.Kind == ast.InputObject) {
			continue
		}
		if strings.HasPrefix(def.Name, "__") {
			continue
		}
		if q, m := b.schema.Query, b.schema.Mutation; (q != nil && def.Name == q.Name) || (m != nil && def.Name == m.Name) {
			continue
		}
		proto, err := gqlutil.ExtractProtoDirective(def.Directives)
		if err != nil {
			return errors.Wrapf(err, "%s has invalid directive", def.Name)
		}
		if proto != nil {
			continue
		}
		if ok, err := b.HasProto(def); err != nil {
			return errors.WithStack(err)
		} else if !ok {
			continue
		}

		obj, err := b.newObject(def)
		if err != nil {
			return errors.WithStack(err)
		}

		b.ObjectsHasProto = append(b.ObjectsHasProto, obj)
	}
	sort.Slice(b.ObjectsHasProto, func(i, j int) bool { return b.ObjectsHasProto[i].Name < b.ObjectsHasProto[j].Name })

	return nil
}

func (b *Binding) newObject(typ *ast.Definition) (*Object, error) {
	proto, err := gqlutil.ExtractProtoDirective(typ.Directives)
	if err != nil {
		return nil, errors.Wrapf(err, "%s has invalid directive", typ.Name)
	}

	obj := &Object{Name: typ.Name, Proto: proto, TypeDef: typ}
	for _, f := range typ.Fields {
		proto, err := gqlutil.ExtractProtoFieldDirective(f.Directives)
		if err != nil {
			return nil, errors.Wrapf(err, "%s has invalid directive", f.Name)
		}

		def := b.schema.Types[f.Type.Name()]

		obj.Fields = append(obj.Fields, &Field{Name: f.Name, GQL: f, Proto: proto, List: f.Type.NamedType == "", TypeDef: def, Object: obj})
	}

	return obj, nil
}

type Binding struct {
	schema          *ast.Schema
	data            *codegen.Data
	ProtoObjects    []*Object
	ObjectsHasProto []*Object
	Enums           []*Enum
}

func (b *Binding) FindGQLFieldType(f *Field) (string, error) {
	if f.Proto == nil {
		if b.data != nil {
			byName := b.data.Objects.ByName
			if f.Object.TypeDef.Kind == ast.InputObject {
				byName = b.data.Inputs.ByName
			}
			obj := byName(f.Object.Name)
			for _, field := range obj.Fields {
				if field.Name == f.Name {
					t := field.TypeReference.GO
					for {
						p, ok := t.(interface{ Elem() types.Type })
						if !ok {
							break
						}
						t = p.Elem()
					}
					return templates.CurrentImports.LookupType(t), nil
				}
			}
		}
		return f.GQL.Type.Name(), nil
	}

	// FIXME
	if f.IsBuiltinType() {
		return f.Proto.Type, nil
	}
	switch f.Proto.Type {
	case "google.protobuf.Int32Value", "google.protobuf.Int64Value",
		"google.protobuf.UInt32Value", "google.protobuf.UInt64Value":
		return "Int", nil
	case "google.protobuf.FloatValue", "google.protobuf.DoubleValue":
		return "Int", nil
	case "google.protobuf.BoolValue":
		return "Boolean", nil
	case "google.protobuf.StringValue":
		return "String", nil
	case "google.protobuf.Timestamp":
		return "DateTime", nil
	}
	for _, o := range b.ProtoObjects {
		if o.Proto.FullName == f.Proto.Type {
			return o.Name, nil
		}
	}
	for _, e := range b.Enums {
		if e.Proto.FullName == f.Proto.Type {
			return e.Name, nil
		}
	}
	return "", errors.Errorf("corresponding GraphQL type was not found: %s", f.Proto.Type)
}

func (b *Binding) HasProto(def *ast.Definition) (bool, error) {
	return gqlutil.HasProto(def, b.schema.Types)
}

type Object struct {
	Name    string
	Proto   *gqlutil.ProtoDirective
	Fields  []*Field
	TypeDef *ast.Definition
}

type Field struct {
	Name    string
	GQL     *ast.FieldDefinition
	Proto   *gqlutil.ProtoFieldDirective
	List    bool
	TypeDef *ast.Definition
	Object  *Object
}

func (f *Field) IsWrapperType() bool {
	if f.Proto == nil {
		return false
	}
	switch f.Proto.Type {
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

func (f *Field) IsBuiltinType() bool {
	if f.Proto == nil {
		// FIXME
		switch f.GQL.Type.Name() {
		case "ID", "Int", "Float", "String", "Boolean":
			return true
		default:
			return false
		}
	}
	return strings.ToLower(f.Proto.Type) == f.Proto.Type
}

func (f *Field) TypeProto() *gqlutil.ProtoDirective {
	if f.TypeDef == nil {
		return nil
	}
	proto, err := gqlutil.ExtractProtoDirective(f.TypeDef.Directives)
	if err != nil {
		return nil
	}
	return proto
}

type Enum struct {
	Name   string
	Proto  *gqlutil.ProtoDirective
	Values []*EnumValue
}

type EnumValue struct {
	Name string
}

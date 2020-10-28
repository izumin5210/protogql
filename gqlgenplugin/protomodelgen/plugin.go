package protomodelgen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/gqlutil"
)

type Plugin struct {
}

func New() *Plugin {
	return new(Plugin)
}

var (
	_ plugin.ConfigMutator
)

func (p *Plugin) Name() string { return "protomodelgen" }

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	binding, err := createBinding(cfg.Schema)
	if err != nil {
		return err
	}

	for _, obj := range binding.Objects {
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

	for _, name := range []string{"Int", "ID"} {
		model := cfg.Models[name]
		model.Model = append(model.Model, "github.com/izumin5210/remixer/gqlruntime.Uint64")
		cfg.Models[name] = model
	}

	return templates.Render(templates.Options{
		PackageName:     cfg.Model.Package,
		Filename:        filepath.Join(cfg.Model.Dir(), "protomodels_gen.go"),
		Data:            binding,
		GeneratedHeader: true,
		Packages:        cfg.Packages,
		Funcs: template.FuncMap{
			"findGQLFieldType": binding.FindGQLFieldType,
		},
	})
}

type ProtoField struct {
	FullName  string
	Package   string
	Name      string
	GoPackage string
	GoName    string
}

func createBinding(s *ast.Schema) (*Binding, error) {
	binding := new(Binding)

	for _, typ := range s.Types {
		proto, err := gqlutil.ExtractProtoDirective(typ.Directives)
		if err != nil {
			return nil, err
		}
		if proto == nil {
			continue
		}

		switch typ.Kind {
		case ast.Object, ast.InputObject:
			obj := &Object{Name: typ.Name, Proto: proto}
			for _, f := range typ.Fields {
				proto, err := gqlutil.ExtractProtoFieldDirective(f.Directives)
				if err != nil {
					return nil, err
				}
				obj.Fields = append(obj.Fields, &Field{Name: f.Name, GQL: f, Proto: proto, List: f.Type.NamedType == ""})
			}
			binding.Objects = append(binding.Objects, obj)

		case ast.Enum:
			enum := &Enum{Name: typ.Name, Proto: proto}
			for _, ev := range typ.EnumValues {
				enum.Values = append(enum.Values, &EnumValue{Name: ev.Name})
			}
			binding.Enums = append(binding.Enums, enum)
		}
	}

	sort.Slice(binding.Objects, func(i, j int) bool { return binding.Objects[i].Name < binding.Objects[j].Name })
	sort.Slice(binding.Enums, func(i, j int) bool { return binding.Enums[i].Name < binding.Enums[j].Name })

	return binding, nil
}

type Binding struct {
	Objects []*Object
	Enums   []*Enum
}

func (b *Binding) FindGQLFieldType(f *Field) (string, error) {
	if f.Proto == nil {
		return f.GQL.Type.Name(), nil
	}

	// FIXME
	if f.IsBuiltinType() {
		return f.Proto.Type, nil
	}
	for _, o := range b.Objects {
		if o.Proto.FullName == f.Proto.Type {
			return o.Name, nil
		}
	}
	for _, e := range b.Enums {
		if e.Proto.FullName == f.Proto.Type {
			return e.Name, nil
		}
	}
	return "", fmt.Errorf("corresponding GraphQL type was not found: %s", f.Proto.Type)
}

type Object struct {
	Name   string
	Proto  *gqlutil.ProtoDirective
	Fields []*Field
}

type Field struct {
	Name  string
	GQL   *ast.FieldDefinition
	Proto *gqlutil.ProtoFieldDirective
	List  bool
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

type Enum struct {
	Name   string
	Proto  *gqlutil.ProtoDirective
	Values []*EnumValue
}

type EnumValue struct {
	Name string
}

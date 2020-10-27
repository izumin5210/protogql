package protomodelgen

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
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
	}
	for _, enum := range binding.Enums {
		cfg.Models.Add(enum.Name, cfg.Model.ImportPath()+"."+enum.Name)
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
		proto, err := extractProtoDirective(typ.Directives)
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
				proto, err := extractProtoFieldDirective(f.Directives)
				if err != nil {
					return nil, err
				}
				if proto == nil {
					continue
				}
				obj.Fields = append(obj.Fields, &Field{Name: f.Name, Proto: proto, List: f.Type.NamedType == ""})
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
	// FIXME
	if strings.ToLower(f.Proto.Type) == f.Proto.Type {
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
	Proto  *ProtoDirective
	Fields []*Field
}

type Field struct {
	Name  string
	Proto *ProtoFieldDirective
	List  bool
}

func (f *Field) GoType() string {
	var sb strings.Builder
	if f.List {
		sb.WriteString("[]")
	}
	sb.WriteString("*")
	// ...
	sb.WriteString(f.Proto.GoType)
	return sb.String()
}

type Enum struct {
	Name   string
	Proto  *ProtoDirective
	Values []*EnumValue
}

type EnumValue struct {
	Name string
}

func extractProtoDirective(directives ast.DirectiveList) (*ProtoDirective, error) {
	for _, d := range directives {
		if d.Name == "proto" {
			out := new(ProtoDirective)
			for _, arg := range d.Arguments {
				switch {
				case arg.Name == "fullName" && arg.Value.Kind == ast.StringValue:
					out.FullName = arg.Value.Raw
				case arg.Name == "package" && arg.Value.Kind == ast.StringValue:
					out.Package = arg.Value.Raw
				case arg.Name == "name" && arg.Value.Kind == ast.StringValue:
					out.Name = arg.Value.Raw
				case arg.Name == "goPackage" && arg.Value.Kind == ast.StringValue:
					out.GoPackage = arg.Value.Raw
				case arg.Name == "goName" && arg.Value.Kind == ast.StringValue:
					out.GoName = arg.Value.Raw
				}
			}
			if !out.IsValid() {
				return nil, fmt.Errorf("invalid proto directive: %w", ErrInvalidArgument)
			}
			return out, nil
		}
	}
	return nil, nil
}

func extractProtoFieldDirective(directives ast.DirectiveList) (*ProtoFieldDirective, error) {
	for _, d := range directives {
		if d.Name == "protoField" {
			out := new(ProtoFieldDirective)
			for _, arg := range d.Arguments {
				switch {
				case arg.Name == "name" && arg.Value.Kind == ast.StringValue:
					out.Name = arg.Value.Raw
				case arg.Name == "type" && arg.Value.Kind == ast.StringValue:
					out.Type = arg.Value.Raw
				case arg.Name == "goName" && arg.Value.Kind == ast.StringValue:
					out.GoName = arg.Value.Raw
				case arg.Name == "goType" && arg.Value.Kind == ast.StringValue:
					out.GoType = arg.Value.Raw
				}
			}
			if !out.IsValid() {
				return nil, fmt.Errorf("invalid protoField directive: %w", ErrInvalidArgument)
			}
			return out, nil
		}
	}
	return nil, nil
}

type ProtoDirective struct {
	FullName  string
	Package   string
	Name      string
	GoPackage string
	GoName    string
}

var ErrInvalidArgument = errors.New("invalid argument")

func (d *ProtoDirective) IsValid() bool {
	return d.FullName != "" && d.Package != "" && d.Name != "" && d.GoPackage != "" && d.GoName != ""
}

type ProtoFieldDirective struct {
	Name   string
	Type   string
	GoName string
	GoType string
}

func (d *ProtoFieldDirective) IsValid() bool {
	return d.Name != "" && d.Type != "" && d.GoName != "" && d.GoType != ""
}

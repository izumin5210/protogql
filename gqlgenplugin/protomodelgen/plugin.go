package protomodelgen

import (
	"path/filepath"
	"text/template"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/pkg/errors"

	"github.com/izumin5210/protogql/gqlgenplugin/protomodelgen/types"
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
		"goWrapperTypeName": types.GoWrapperTypeName,
		"unwrapStatement":   types.UnwrapStatement,
		"isProtoType":       types.IsProtoType,
	}
)

func (p *Plugin) Name() string { return "protomodelgen" }

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	reg, err := types.CreateRegistryFromSchema(cfg.Schema)
	if err != nil {
		return err
	}

	models := []types.ProtoType{}
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
			if f.IsDefinedInProto() {
				continue
			}
			if cfg.Models[obj.GQLName()].Fields == nil {
				cfg.Models[obj.GQLName()] = config.TypeMapEntry{
					Model:  cfg.Models[obj.GQLName()].Model,
					Fields: map[string]config.TypeMapField{},
				}
			}
			cfg.Models[obj.GQLName()].Fields[f.GQLName()] = config.TypeMapField{FieldName: f.GQLName(), Resolver: true}
		}
	}

	for _, enum := range reg.EnumsFromProto() {
		cfg.Models.Add(enum.GQLName(), cfg.Model.ImportPath()+"."+enum.GoTypeName())
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
	reg, err := types.CreateRegistry(data)
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

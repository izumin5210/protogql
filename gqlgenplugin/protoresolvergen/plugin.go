package protoresolvergen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/tools/go/packages"

	"github.com/izumin5210/remixer/goutil"
	"github.com/izumin5210/remixer/gqlutil"
)

type Plugin struct {
}

func New() *Plugin {
	return new(Plugin)
}

var (
	_ plugin.CodeGenerator
)

func (p *Plugin) Name() string { return "protomodelgen" }

func (p *Plugin) GenerateCode(data *codegen.Data) error {
	if !data.Config.Resolver.IsDefined() {
		return nil
	}

	switch data.Config.Resolver.Layout {
	case config.LayoutSingleFile:
		return p.generateSingleFile(data)
	case config.LayoutFollowSchema:
		return p.generatePerSchema(data)
	}

	return nil
}

func (p *Plugin) generateSingleFile(data *codegen.Data) error {
	return nil
}

func (p *Plugin) generatePerSchema(data *codegen.Data) error {
	files := NewFiles(data.Config.Resolver)

	// 生成したいファイル
	// * .graphql ファイルに対応する .proto.resolvers.go
	// * Proto Resolver -> Model Resolver の変換をする adapter 実装がまとまったファイル

	for _, o := range data.Objects {
		if !o.HasResolvers() {
			continue
		}
		file := files.FindOrInitialize(o.Position.Src.Name)
		file.Objects = append(file.Objects, o)

		for _, field := range o.Fields {
			if !field.IsResolver {
				continue
			}

			gqlType := data.Config.Schema.Types[field.TypeReference.GQL.Name()]

			var modelPkg *packages.Package
			if typeRef := field.TypeReference; !typeRef.IsScalar() {
				modelPkg = data.Config.Packages.LoadWithTypes(goutil.GetTypePackageName(typeRef.GO))
			}

			file := files.FindOrInitialize(field.Position.Src.Name)
			file.Resolvers = append(file.Resolvers, &Resolver{Field: field, GQLTypeDefinition: gqlType, modelPkg: modelPkg})
		}
	}

	for filename, file := range files.files {
		err := templates.Render(templates.Options{
			PackageName: data.Config.Resolver.Package,
			Filename:    filename,
			Data:        file,
			Funcs: template.FuncMap{
				"resolverImplementationName": file.ResolverImplementationName,
				"resolverAdapterName":        file.ResolverAdapterName,
			},
			Packages: data.Config.Packages,
		})

		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(data.Config.Resolver.Filename); errors.Is(err, os.ErrNotExist) {
		err := templates.Render(templates.Options{
			PackageName: data.Config.Resolver.Package,
			FileNotice: `
				// This file will not be regenerated automatically.
				//
				// It serves as dependency injection for your app, add any dependencies you require here.`,
			Template: `type {{.}} struct {}`,
			Filename: data.Config.Resolver.Filename,
			Data:     data.Config.Resolver.Type,
			Packages: data.Config.Packages,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type Files struct {
	files map[string]*File
	cfg   config.ResolverConfig
}

func NewFiles(cfg config.ResolverConfig) *Files {
	return &Files{files: map[string]*File{}, cfg: cfg}
}

func (f *Files) FindOrInitialize(name string) *File {
	filename := f.resolverGoFilename(name)
	if _, ok := f.files[filename]; !ok {
		f.files[filename] = &File{ResolverType: f.cfg.Type}
	}
	return f.files[filename]
}

func (f *Files) resolverGoFilename(gqlFilename string) string {
	tmpl := f.cfg.FilenameTemplate
	if tmpl == "" {
		tmpl = "{name}.resolvers.proto.go"
	} else {
		tmpl = strings.TrimSuffix(tmpl, ".go") + ".proto.go"
	}
	filename := strings.TrimSuffix(filepath.Base(gqlFilename), filepath.Ext(gqlFilename))
	filename = strings.ReplaceAll(tmpl, "{name}", filename)
	return filepath.Join(f.cfg.Dir(), filename)
}

type Resolver struct {
	*codegen.Field
	GQLTypeDefinition *ast.Definition
	modelPkg          *packages.Package
}

type ResolverArg struct {
	Name, Type string
}

func (r *Resolver) ShortProtoResolverDeclaration() (string, error) {
	var proto *gqlutil.ProtoDirective
	var err error

	proto, err = gqlutil.ExtractProtoDirective(r.GQLTypeDefinition.Directives)
	if err != nil {
		return "", err
	}

	args, err := r.ResolverArgs()
	if err != nil {
		return "", err
	}
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = fmt.Sprintf("%s %s", arg.Name, arg.Type)
	}

	var result string
	if proto == nil {
		result = templates.CurrentImports.LookupType(r.TypeReference.GO)
	} else {
		result = fmt.Sprintf("*%s.%s", templates.CurrentImports.Lookup(proto.GoPackage), proto.GoName)
		if r.Type.Elem != nil {
			result = "[]" + result
		}
	}

	return fmt.Sprintf("(%s) (%s, error)", strings.Join(argStrs, ", "), result), nil
}

func (r *Resolver) ResolverArgs() ([]ResolverArg, error) {
	args := []ResolverArg{{Name: "ctx", Type: "context.Context"}}

	if !r.Object.Root {
		parentProto, err := gqlutil.ExtractProtoDirective(r.Object.Definition.Directives)
		if err != nil {
			return nil, err
		}
		var typ string
		if parentProto == nil {
			typ = templates.CurrentImports.LookupType(r.Object.Reference())
		} else {
			typ = fmt.Sprintf("*%s.%s", templates.CurrentImports.Lookup(parentProto.GoPackage), parentProto.GoName)
		}
		args = append(args, ResolverArg{Name: "obj", Type: typ})
	}

	for _, arg := range r.Args {
		// TODO: proto にする
		args = append(args, ResolverArg{Name: arg.VarName, Type: templates.CurrentImports.LookupType(arg.TypeReference.GO)})
	}

	return args, nil
}

func (r *Resolver) ArgList() (string, error) {
	argNames := []string{}

	args, err := r.ResolverArgs()
	if err != nil {
		return "", nil
	}

	for _, arg := range args {
		n := arg.Name
		if n == "obj" {
			n = r.ModelToProtoFunc(r.Object.Definition, false) + "(obj)"
		}
		argNames = append(argNames, n)
	}

	return strings.Join(argNames, ", "), nil
}

func (r *Resolver) ResolverModelFromProtoFunc() string {
	return r.ModelFromProtoFunc(r.TypeReference.Definition, r.TypeReference.IsSlice())
}

func (r *Resolver) ModelFromProtoFunc(def *ast.Definition, slice bool) string {
	return r.modelMappingFunc(def, slice, "From")
}

func (r *Resolver) ModelToProtoFunc(def *ast.Definition, slice bool) string {
	return r.modelMappingFunc(def, slice, "To")
}

func (r *Resolver) modelMappingFunc(def *ast.Definition, slice bool, d string) string {
	var list, suffix string

	if slice {
		list = "List"
		suffix = "RepeatedProto"
	} else {
		suffix = "Proto"
	}

	return fmt.Sprintf("%s.%s%s%s%s", templates.CurrentImports.Lookup(r.modelPkg.PkgPath), def.Name, list, d, suffix)
}

func (r *Resolver) IsList() bool { return r.Type.Elem != nil }

type File struct {
	Objects      []*codegen.Object
	Resolvers    []*Resolver
	ResolverType string
}

func (f *File) ResolverImplementationName(obj *codegen.Object) string {
	return fmt.Sprintf("%sProto%s", templates.LcFirst(obj.Name), templates.UcFirst(f.ResolverType))
}

func (f *File) ResolverAdapterName(obj *codegen.Object) string {
	return fmt.Sprintf("%sProto%sAdapter", templates.LcFirst(obj.Name), templates.UcFirst(f.ResolverType))
}

package protoresolvergen

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/tools/go/packages"

	"github.com/izumin5210/protogql/codegen/goutil"
	"github.com/izumin5210/protogql/codegen/gqlutil"
)

type Plugin struct {
}

func New() *Plugin {
	return new(Plugin)
}

var (
	_ plugin.CodeGenerator
)

func (p *Plugin) Name() string { return "protoresolvergen" }

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
	files := NewFiles(data.Config.Resolver, data.Config.Model, data.Schema)

	var (
		pkg *packages.Package

		// ref: https://github.com/99designs/gqlgen/blob/v0.13.0/internal/rewrite/rewriter.go
		copiedDecls = map[goast.Decl]struct{}{}
	)

	if importPath := goutil.GetImportPathForDir(data.Config.Resolver.Dir()); importPath != "" {
		pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedSyntax | packages.NeedTypes}, importPath)
		if err == nil && len(pkgs) > 0 {
			pkg = pkgs[0]
		}
	}

	for _, o := range data.Objects {
		if !o.HasResolvers() {
			continue
		}
		file := files.FindOrInitialize(o.Position.Src.Name)
		file.Objects = append(file.Objects, o)

		if pkg != nil {
			if d := goutil.GetStruct(pkg, file.ProtoResolverImplementationName(o)); d != nil {
				copiedDecls[d] = struct{}{}
			}
			if d := goutil.GetMethod(pkg, file.ResolverTypeName(), o.Name); d != nil {
				copiedDecls[d] = struct{}{}
			}
		}

		for _, field := range o.Fields {
			if !field.IsResolver {
				continue
			}

			gqlType := data.Config.Schema.Types[field.TypeReference.GQL.Name()]

			var modelPkg *packages.Package
			if typeRef := field.TypeReference; !typeRef.IsScalar() {
				modelPkg = data.Config.Packages.LoadWithTypes(goutil.GetTypePackageName(typeRef.GO))
			}
			protoField, err := gqlutil.ExtractProtoFieldDirective(field.FieldDefinition.Directives)
			if err != nil {
				return errors.Wrapf(err, "%s has invalid directive", field.Name)
			}

			file := files.FindOrInitialize(field.Position.Src.Name)
			r := &Resolver{Field: field, ProtoField: protoField, GQLTypeDefinition: gqlType, modelPkg: modelPkg, cfg: data.Config.Resolver, file: file, pkg: pkg}
			file.Resolvers = append(file.Resolvers, r)

			if pkg != nil {
				if d := goutil.GetMethod(pkg, r.ProtoImplementationName(), r.GoFieldName); d != nil {
					copiedDecls[d] = struct{}{}
				}
			}
		}
	}

	for _, file := range files.Files {
		if pkg != nil {
			if goFile := goutil.GetFile(pkg, file.ProtoResolverGoFilename()); goFile != nil {
				var buf bytes.Buffer
				for _, d := range goFile.Decls {
					if _, ok := copiedDecls[d]; ok {
						continue
					}
					if gd, ok := d.(*goast.GenDecl); ok && gd.Tok == token.IMPORT {
						continue
					}
					source, err := goutil.GetSource(pkg.Fset, d.Pos(), d.End())
					if err != nil {
						return err
					}
					buf.WriteString(source)
					buf.WriteString("\n")
				}
				file.ProtoResolverRemainingSource = buf.String()
				file.Imports = goutil.ListImports(goFile)
			}
		}

		err := templates.Render(templates.Options{
			PackageName: data.Config.Resolver.Package,
			Template:    templateProtoResolvers,
			Filename:    file.ProtoResolverGoFilename(),
			Data:        file,
			Packages:    data.Config.Packages,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to render %s", file.ProtoResolverGoFilename())
		}

		err = templates.Render(templates.Options{
			PackageName: data.Config.Resolver.Package,
			Template:    templateResolvers,
			Filename:    file.ResolverGoFilename(),
			Data:        file,
			Packages:    data.Config.Packages,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to render %s", file.ResolverGoFilename())
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
			return errors.Wrapf(err, "failed to render %s", data.Config.Resolver.Filename)
		}
	}
	return nil
}

type Files struct {
	Files    map[string]*File
	cfg      config.ResolverConfig
	modelCfg config.PackageConfig
	schema   *ast.Schema
}

func NewFiles(cfg config.ResolverConfig, modelCfg config.PackageConfig, schema *ast.Schema) *Files {
	return &Files{Files: map[string]*File{}, cfg: cfg, modelCfg: modelCfg, schema: schema}
}

func (f *Files) FindOrInitialize(gqlFilename string) *File {
	if _, ok := f.Files[gqlFilename]; !ok {
		f.Files[gqlFilename] = &File{gqlFilename: gqlFilename, cfg: f.cfg, modelCfg: f.modelCfg, schema: f.schema}
	}
	return f.Files[gqlFilename]
}

type Resolver struct {
	*codegen.Field
	file              *File
	cfg               config.ResolverConfig
	ProtoField        *gqlutil.ProtoFieldDirective
	GQLTypeDefinition *ast.Definition
	modelPkg          *packages.Package
	pkg               *packages.Package
}

type ResolverArg struct {
	Name, Type string
}

func (r *Resolver) ShortProtoResolverDeclaration() (string, error) {
	var proto *gqlutil.ProtoDirective
	var err error

	proto, err = gqlutil.ExtractProtoDirective(r.GQLTypeDefinition.Directives)
	if err != nil {
		return "", errors.Wrapf(err, "%s has invalid directive", r.GQLTypeDefinition.Name)
	}

	args, err := r.ResolverArgs()
	if err != nil {
		return "", errors.WithStack(err)
	}
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = fmt.Sprintf("%s %s", arg.Name, arg.Type)
	}

	var result string
	if proto == nil {
		typeDef := r.GQLTypeDefinition
		if typeDef.Kind == ast.Object || typeDef.Kind == ast.InputObject {
			if ok, err := gqlutil.HasProto(typeDef, r.file.schema.Types); err == nil && ok {
				result = fmt.Sprintf("*%s.%s", templates.CurrentImports.Lookup(r.file.modelCfg.ImportPath()), typeDef.Name+"_Proto")
			}
		}
		if result == "" {
			result = templates.CurrentImports.LookupType(r.TypeReference.GO)
		}
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
			return nil, errors.Wrapf(err, "%s has invalid directive", r.Object.Name)
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

func (r *Resolver) ProtoImplementationName() string {
	return r.file.ProtoResolverImplementationName(r.Object)
}

func (r *Resolver) ImplementationName() string {
	return r.file.ResolverImplementationName(r.Object)
}

func (r *Resolver) ProtoResolverBody() string {
	notImplemented := `panic("not implemented")`
	if r.pkg == nil {
		return notImplemented
	}
	meth := goutil.GetMethod(r.pkg, r.ProtoImplementationName(), r.GoFieldName)
	if meth == nil {
		return notImplemented
	}

	body, err := goutil.GetSource(r.pkg.Fset, meth.Body.Pos()+1, meth.Body.End()-1) // Avoid braces
	if err == nil && body != "" {
		return strings.TrimSpace(body)
	}

	return notImplemented
}

type File struct {
	gqlFilename                  string
	cfg                          config.ResolverConfig
	modelCfg                     config.PackageConfig
	schema                       *ast.Schema
	Objects                      []*codegen.Object
	Resolvers                    []*Resolver
	Imports                      []*goutil.Import
	ProtoResolverRemainingSource string
}

func (f *File) ResolverTypeName() string {
	return f.cfg.Type
}

// https://github.com/99designs/gqlgen/blob/v0.13.0/plugin/resolvergen/resolver.go#L199-L207
func (f *File) ResolverGoFilename() string {
	tmpl := f.cfg.FilenameTemplate
	if tmpl == "" {
		tmpl = "{name}.resolvers.go"
	}
	filename := strings.TrimSuffix(filepath.Base(f.gqlFilename), filepath.Ext(f.gqlFilename))
	filename = strings.ReplaceAll(tmpl, "{name}", filename)
	return filepath.Join(f.cfg.Dir(), filename)
}

func (f *File) ProtoResolverGoFilename() string {
	return strings.TrimSuffix(f.ResolverGoFilename(), ".go") + ".proto.go"
}

func (f *File) ProtoResolverImplementationName(obj *codegen.Object) string {
	return fmt.Sprintf("%sProto%s", templates.LcFirst(obj.Name), templates.UcFirst(f.cfg.Type))
}

func (f *File) ResolverImplementationName(obj *codegen.Object) string {
	return fmt.Sprintf("%s%s", templates.LcFirst(obj.Name), templates.UcFirst(f.cfg.Type))
}

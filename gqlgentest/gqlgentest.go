package gqlgentest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type Runner interface {
	Run(t *testing.T, f func(t *testing.T, err error))
	AddGqlGenOption(opst ...api.Option)
	AddGqlSchema(filename, content string)
	AddGqlSchemaFile(t *testing.T, pattern string)
	AddGoModReplace(pkg, path string)
	Snapshot(t *testing.T, v ...interface{})
	SnapshotFile(t *testing.T, files ...string)
}

func New(t *testing.T) Runner {
	dir := filepath.Join(t.TempDir(), "testapp")
	gqlgenCfg := config.DefaultConfig()
	gqlgenCfg.Exec = config.PackageConfig{
		Filename: "graph/graph_gen.go",
		Package:  "graph",
	}
	gqlgenCfg.Model = config.PackageConfig{
		Filename: "model/models_gen.go",
		Package:  "model",
	}
	gqlgenCfg.Resolver = config.ResolverConfig{
		Layout:  config.LayoutFollowSchema,
		DirName: "resolver",
		Package: "resolver",
	}
	return &runner{dir: dir, gqlgenCfg: gqlgenCfg}
}

type runner struct {
	dir           string
	prevDir       string
	gqlgenCfg     *config.Config
	gqlgenOptions []api.Option
	goModReplace  []struct{ Package, Path string }
}

func (r *runner) Run(t *testing.T, f func(t *testing.T, err error)) {
	defer r.moveToRoot(t)()
	r.writeGoMod(t)

	r.orDie(t, api.Generate(r.gqlgenCfg, r.gqlgenOptions...))

	f(t, nil)
}

func (r *runner) orDie(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func (r *runner) moveToRoot(t *testing.T) func() {
	if r.prevDir != "" {
		t.Fatal("already moved to testapp root")
	}

	var err error

	r.prevDir, err = os.Getwd()
	r.orDie(t, err)

	err = os.Mkdir(r.dir, os.FileMode(0755))
	r.orDie(t, err)

	err = os.Chdir(r.dir)
	r.orDie(t, err)

	return func() {
		err := os.Chdir(r.prevDir)
		r.orDie(t, err)
	}
}

func (r *runner) AddGqlGenOption(opts ...api.Option) {
	r.gqlgenOptions = append(r.gqlgenOptions, opts...)
}

func (r *runner) AddGqlSchema(filename, content string) {
	r.gqlgenCfg.Sources = append(r.gqlgenCfg.Sources, &ast.Source{Name: filename, Input: content})
}

func (r *runner) AddGqlSchemaFile(t *testing.T, pattern string) {
	paths, err := filepath.Glob(pattern)
	r.orDie(t, err)

	for _, path := range paths {
		st, err := os.Stat(path)
		r.orDie(t, err)
		if st.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(path)
		r.orDie(t, err)
		r.AddGqlSchema(path, string(data))
	}
}

func (r *runner) AddGoModReplace(pkg, path string) {
	r.goModReplace = append(r.goModReplace, struct{ Package, Path string }{pkg, path})
}

func (r *runner) writeGoMod(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("module testapp\n")
	buf.WriteString("go 1.15\n")

	for _, r := range r.goModReplace {
		buf.WriteString(fmt.Sprintf("replace %s => %s\n", r.Package, r.Path))
	}

	r.orDie(t, ioutil.WriteFile(filepath.Join(r.dir, "go.mod"), buf.Bytes(), os.FileMode(0644)))
}

func (r *runner) Snapshot(t *testing.T, v ...interface{}) {
	t.Helper()
	cupaloy.Global.WithOptions(cupaloy.SnapshotSubdirectory(filepath.Join(r.prevDir, ".snapshots"))).SnapshotT(t, v...)
}

func (r *runner) SnapshotFile(t *testing.T, files ...string) {
	t.Helper()
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				t.Errorf("failed to read file: %v", err)
			}
			r.Snapshot(t, string(data))
		})
	}
}

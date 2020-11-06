package goutil

import (
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func GetTypePackageName(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return GetTypePackageName(t.Elem())
	case *types.Slice:
		return GetTypePackageName(t.Elem())
	case *types.Named:
		return t.Obj().Pkg().Path()
	default:
		panic(errors.Errorf("%T is not supported", t))
	}
}

func GetFile(pkg *packages.Package, filename string) *ast.File {
	for _, gf := range pkg.Syntax {
		if pkg.Fset.Position(gf.Pos()).Filename == filename {
			return gf
		}
	}
	return nil
}

type Import struct{ Alias, Path string }

func ListImports(f *ast.File) []*Import {
	imports := make([]*Import, len(f.Imports))

	for i, imp := range f.Imports {
		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		}
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			panic(err)
		}
		imports[i] = &Import{Alias: name, Path: path}
	}

	return imports
}

func GetStruct(pkg *packages.Package, structName string) *ast.GenDecl {
	for _, s := range pkg.Syntax {
		for _, d := range s.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gd.Tok != token.TYPE {
				continue
			}

			for _, s := range gd.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if ts.Name.Name != structName {
					continue
				}

				return gd
			}
		}
	}

	return nil
}

func GetMethod(pkg *packages.Package, structName, methodName string) *ast.FuncDecl {
	for _, s := range pkg.Syntax {
		for _, d := range s.Decls {
			m, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if m.Recv == nil || len(m.Recv.List) == 0 {
				continue
			}
			if m.Name.Name != methodName {
				continue
			}

			recvType := m.Recv.List[0].Type
			if s, ok := recvType.(*ast.StarExpr); ok {
				recvType = s.X
			}
			if ident, ok := recvType.(*ast.Ident); !ok || ident.Name != structName {
				continue
			}

			return m
		}
	}

	return nil
}

func GetSource(fset *token.FileSet, start, end token.Pos) (string, error) {
	startPos, endPos := fset.Position(start), fset.Position(end)
	data, err := ioutil.ReadFile(startPos.Filename)
	if err != nil {
		return "", errors.Wrap(err, "failed to read file")
	}

	return string(data[startPos.Offset:endPos.Offset]), nil
}

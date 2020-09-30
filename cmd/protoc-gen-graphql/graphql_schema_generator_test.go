package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoprocessor"
)

func TestProcessor(t *testing.T) {
	testGenerate(t, "user")
	testGenerate(t, "starwars")
}

func getFixtureTypes(t *testing.T, protosetName string) *protoprocessor.Types {
	f, err := ioutil.ReadFile(filepath.Join("testdata", protosetName))
	if err != nil {
		t.Fatalf("failed to open fixture: %v", err)
	}

	var set descriptor.FileDescriptorSet
	err = proto.Unmarshal(f, &set)
	if err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	types := protoprocessor.NewTypes()
	for _, f := range set.GetFile() {
		types.AddFile(f)
	}

	return types
}

func pickFile(name string, files []*plugin.CodeGeneratorResponse_File) (picked *plugin.CodeGeneratorResponse_File, rest []*plugin.CodeGeneratorResponse_File) {
	for _, f := range files {
		if f.GetName() == name {
			picked = f
		} else {
			rest = append(rest, f)
		}
	}
	return
}

func testGenerate(t *testing.T, fixture string) {
	t.Run(fixture, func(t *testing.T) {
		types := getFixtureTypes(t, fixture+".protoset")

		schemata := []*ast.Source{{Input: BaseSchema}}

		for _, filename := range types.Files() {
			filename := filename
			name := filename[strings.LastIndex(filename, "/")+1:]

			t.Run(name, func(t *testing.T) {
				file, err := GraphQLSchemaGenerator.Generate(context.Background(), filename, types)
				if err != nil {
					t.Errorf("Generate() returns %v, want nil", err)
				}

				if file == nil {
					return
				}

				cupaloy.SnapshotT(t, file.GetContent())
				schemata = append(schemata, &ast.Source{Input: file.GetContent()})
			})
		}

		_, gqlErr := gqlparser.LoadSchema(schemata...)
		if gqlErr != nil {
			t.Errorf("generated schema has some violations:\n%v", gqlErr)
		}
	})
}

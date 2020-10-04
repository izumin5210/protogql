package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestProcessor(t *testing.T) {
	testGenerate(t, "user")
	testGenerate(t, "starwars")
}

func getFixtureFiles(t *testing.T, protosetName string) *protoregistry.Files {
	f, err := ioutil.ReadFile(filepath.Join("testdata", protosetName))
	if err != nil {
		t.Fatalf("failed to open fixture: %v", err)
	}

	var set descriptorpb.FileDescriptorSet
	err = proto.Unmarshal(f, &set)
	if err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	files, err := protodesc.NewFiles(&set)
	if err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	return files
}

func pickFile(name string, files []*pluginpb.CodeGeneratorResponse_File) (picked *pluginpb.CodeGeneratorResponse_File, rest []*pluginpb.CodeGeneratorResponse_File) {
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
		files := getFixtureFiles(t, fixture+".protoset")

		schemata := []*ast.Source{{Input: BaseSchema}}

		files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
			if !strings.HasPrefix(fd.Path(), filepath.Join("testdata", fixture)) {
				return true
			}

			filename := fd.Path()
			name := filename[strings.LastIndex(filename, "/")+1:]
			name = strings.TrimSuffix(name, ".proto")
			name = name + ".gql"

			t.Run(name, func(t *testing.T) {
				file, err := GraphQLSchemaGenerator.Generate(context.Background(), fd)
				if err != nil {
					t.Errorf("Generate() returns %v, want nil", err)
				}

				if file == nil {
					return
				}

				cupaloy.SnapshotT(t, file.GetContent())
				schemata = append(schemata, &ast.Source{Input: file.GetContent()})
			})

			return true
		})

		_, gqlErr := gqlparser.LoadSchema(schemata...)
		if gqlErr != nil {
			t.Errorf("generated schema has some violations:\n%v", gqlErr)
		}
	})
}

package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
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
	t.Run("user.gql", func(t *testing.T) {
		types := getFixtureTypes(t, "user.protoset")

		file, err := GraphQLSchemaGenerator.Generate(context.Background(), "testdata/user.proto", types)
		if err != nil {
			t.Errorf("Generate() returns %v, want nil", err)
		}

		if file == nil {
			t.Error("user.gql was not generated")
		} else {
			cupaloy.SnapshotT(t, file.GetContent())

			_, gqlErr := gqlparser.LoadSchema(&ast.Source{Name: file.GetName(), Input: file.GetContent()})
			if gqlErr != nil {
				t.Errorf("generated schema has some violations:\n%v", gqlErr)
			}
		}
	})
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

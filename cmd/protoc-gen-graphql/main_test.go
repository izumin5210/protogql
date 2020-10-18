package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestProcessor(t *testing.T) {
	testGenerate(t, "user")
	testGenerate(t, "starwars")
}

func loadDescriptorSet(t *testing.T, protosetName string) *descriptorpb.FileDescriptorSet {
	f, err := ioutil.ReadFile(filepath.Join("testdata", protosetName))
	if err != nil {
		t.Fatalf("failed to open fixture: %v", err)
	}

	var set descriptorpb.FileDescriptorSet
	err = proto.Unmarshal(f, &set)
	if err != nil {
		t.Fatalf("failed to parse fixture: %v", err)
	}

	return &set
}

func newCodeGeneratorRequestFromDescriptorSet(set *descriptorpb.FileDescriptorSet, fileToGeneratePrefix string) *pluginpb.CodeGeneratorRequest {
	req := new(pluginpb.CodeGeneratorRequest)

	for _, f := range set.GetFile() {
		req.ProtoFile = append(req.ProtoFile, f)
		if n := f.GetName(); strings.HasPrefix(n, fileToGeneratePrefix) {
			req.FileToGenerate = append(req.FileToGenerate, n)
		}
	}

	return req
}

func testProtocGen(t *testing.T, descriptorSetPath, fileToGeneratePrefix string) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	req := newCodeGeneratorRequestFromDescriptorSet(
		loadDescriptorSet(t, descriptorSetPath),
		fileToGeneratePrefix,
	)

	plugin, err := options.New(req)
	if err != nil {
		t.Fatalf("faield to create plugin: %v", err)
	}

	err = run(plugin)
	if err != nil {
		t.Errorf("Generator returns %v, want nil", err)
	}

	return plugin.Response()
}

func testGenerate(t *testing.T, fixture string) {
	t.Run(fixture, func(t *testing.T) {
		resp := testProtocGen(t, fixture+".protoset", filepath.Join("testdata", fixture))

		schemata := []*ast.Source{{Input: BaseSchema}}
		for _, f := range resp.GetFile() {
			f := f
			name := f.GetName()
			name = name[strings.LastIndex(name, "/")+1:]
			t.Run(name, func(t *testing.T) {
				cupaloy.SnapshotT(t, f.GetContent())
			})
			schemata = append(schemata, &ast.Source{Input: f.GetContent()})
		}

		_, gqlErr := gqlparser.LoadSchema(schemata...)
		if gqlErr != nil {
			t.Errorf("generated schema has some violations:\n%v", gqlErr)
		}
	})
}

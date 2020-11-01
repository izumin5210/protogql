package plugin_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/izumin5210/remixer/gqlgenplugin/protomodelgen"
	"github.com/izumin5210/remixer/gqlgenplugin/protoresolvergen"
	"github.com/izumin5210/remixer/gqlgentest"
)

func TestGenerateForProto(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenPlugin(protomodelgen.New())
	gqlgentest.AddGqlGenPlugin(protoresolvergen.New())
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "task", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime

extend type Query {
  tasks: [Task!]!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/task", filepath.Join(testdataDir, "apis", "go", "task"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}
		gqlgentest.SnapshotFile(t,
			"resolver/resolver.go",
			"resolver/resolver.adapters.go",
			"resolver/schema.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithExtendingType(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenPlugin(protomodelgen.New())
	gqlgentest.AddGqlGenPlugin(protoresolvergen.New())
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "task", "*.graphqls"))
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "user", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime

extend type Task {
  assignees: [User!]!
  author: User!
}

extend type Query {
  tasks: [Task!]!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/task", filepath.Join(testdataDir, "apis", "go", "task"))
	gqlgentest.AddGoModReplace("apis/go/user", filepath.Join(testdataDir, "apis", "go", "user"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}
		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/resolver.adapters.go",
			"resolver/schema.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithProtoWellKnownTypes(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenPlugin(protomodelgen.New())
	gqlgentest.AddGqlGenPlugin(protoresolvergen.New())
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "wktypes", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime

extend type Query {
  hello: [Hello!]!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/wktypes", filepath.Join(testdataDir, "apis", "go", "wktypes"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}
		gqlgentest.SnapshotFile(t,
			"resolver/resolver.go",
			"resolver/resolver.adapters.go",
			"resolver/schema.resolvers.proto.go",
		)
	})
}

func getModuleRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), ".."))
}

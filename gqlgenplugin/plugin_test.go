package gqlgenplugin_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/99designs/gqlgen/api"
	"github.com/izumin5210/remixer/gqlgenplugin"
	"github.com/izumin5210/remixer/gqlgenplugin/protomodelgen"
	"github.com/izumin5210/remixer/gqlgenplugin/protoresolvergen"
	"github.com/izumin5210/remixer/gqlgentest"
)

func TestGenerateForProto(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.PrependPlugin(protomodelgen.New()),
		api.AddPlugin(protoresolvergen.New()),
		gqlgenplugin.RemovePlugin("resolvergen"),
	)
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
		if entries, err := filepath.Glob("resolver/**"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 3; got != want {
			t.Errorf("Files under resolver/ were found %d, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/schema.resolvers.go",
			"resolver/schema.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithExtendingType(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.PrependPlugin(protomodelgen.New()),
		api.AddPlugin(protoresolvergen.New()),
		gqlgenplugin.RemovePlugin("resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "task", "*.graphqls"))
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "user", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime`)
	gqlgentest.AddGqlSchema("user.graphqls", `
extend type Query {
  currentUser: User!
}`)
	gqlgentest.AddGqlSchema("task.graphqls", `
extend type Task {
  assignees: [User!]!
  author: User!
}

extend type User {
  assignedTasks: [Task!]!
}

extend type Query {
  tasks: [Task!]!
  latestTask: Task!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/task", filepath.Join(testdataDir, "apis", "go", "task"))
	gqlgentest.AddGoModReplace("apis/go/user", filepath.Join(testdataDir, "apis", "go", "user"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}

		if entries, err := filepath.Glob("resolver/**"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 9; got != want {
			t.Errorf("Files under resolver/ were found %d, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/task.pb.resolvers.go",
			"resolver/task.pb.resolvers.proto.go",
			"resolver/task.resolvers.go",
			"resolver/task.resolvers.proto.go",
			"resolver/user.pb.resolvers.go",
			"resolver/user.pb.resolvers.proto.go",
			"resolver/user.resolvers.go",
			"resolver/user.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithProtoWellKnownTypes(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.PrependPlugin(protomodelgen.New()),
		api.AddPlugin(protoresolvergen.New()),
		gqlgenplugin.RemovePlugin("resolvergen"),
	)
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

		if entries, err := filepath.Glob("resolver/**"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 3; got != want {
			t.Errorf("Files under resolver/ were found %d, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/schema.resolvers.go",
			"resolver/schema.resolvers.proto.go",
		)
	})
}

func getModuleRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), ".."))
}

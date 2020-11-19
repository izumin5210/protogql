package gqlgenplugin_test

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

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
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "hello", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION

extend type Query {
  hello: Hello!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/hello", filepath.Join(testdataDir, "apis", "go", "hello"))

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
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "hello", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime`)
	gqlgentest.AddGqlSchema("hello.graphqls", `
extend type Hello {
  user: User!
}

extend type Query {
  hello: Hello!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/hello", filepath.Join(testdataDir, "apis", "go", "hello"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}

		if entries, err := filepath.Glob("resolver/**"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 5; got != want {
			t.Errorf("Files under resolver/ were found %d, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/hello.pb.resolvers.go",
			"resolver/hello.pb.resolvers.proto.go",
			"resolver/hello.resolvers.go",
			"resolver/hello.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithProtoComments(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		// gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "comments", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/comments", filepath.Join(testdataDir, "apis", "go", "comments"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			// "resolver/resolver.go",
			// "resolver/schema.resolvers.go",
			// "resolver/schema.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithProtoWellKnownTypes(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
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

func TestGenerateForProto_WithProtoOneofs(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "oneof", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String, oneofName: String, oneofGoName: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
input ProtoOneof {
  fields: [ProtoOneofField!]!
}
input ProtoOneofField {
  name: String!
  goName: String!
}

extend type Query {
  entries: [Entry!]!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/oneof", filepath.Join(testdataDir, "apis", "go", "oneof"))

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

func TestGenerateForProto_WhenUpdate(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "hello", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime`)
	gqlgentest.AddGqlSchema("hello.graphqls", `
extend type Hello {
  user: User!
  replies: [Hello!]!
}

extend type Query {
  hello: Hello!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/hello", filepath.Join(testdataDir, "apis", "go", "hello"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("failed to generate code: %v", err)
		}

		rewriteFile(t, "resolver/hello.resolvers.proto.go", func(input string, write func(string)) {
			write(input)
			write("const (\n\tTestConstant = 1\n)\n")
			write("var (\n\tTestVariable = 1\n)\n")
			write("type TestStruct struct {\n\tFoo string\n}\n")
			write("func TestFunction() string { return \"Test\" }\n")
		})
		rewriteFile(t, "resolver/hello.resolvers.proto.go", func(input string, write func(string)) {
			lines := strings.Split(input, "\n")
			for i := 0; i < len(lines); i++ {
				write(lines[i])
				write("\n")
				if strings.HasPrefix(lines[i], "func (r *queryProtoResolver) Hello") {
					write("\treturn &hello_pb.Hello{}, nil\n")
					i++
				}
				if lines[i] == "import (" {
					write("\t_ \"net/http/pprof\"\n")
					i++
				}
			}
		})
	})

	gqlgentest.ReplaceGqlSchema("hello.graphqls", `
extend type Hello {
  user: User!
}

extend type Query {
  hello: Hello!
  hellos: [Hello!]!
}`)

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %+v", err)
		}

		if entries, err := filepath.Glob("resolver/**"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 5; got != want {
			t.Errorf("Files under resolver/ were found %d, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/hello.pb.resolvers.go",
			"resolver/hello.pb.resolvers.proto.go",
			"resolver/hello.resolvers.go",
			"resolver/hello.resolvers.proto.go",
		)
	})
}

func TestGenerateForProto_WithPlainGqlTypes(t *testing.T) {
	rootDir := getModuleRoot()
	testdataDir := filepath.Join(rootDir, "testdata")

	gqlgentest := gqlgentest.New(t)
	gqlgentest.AddGqlGenOption(
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	gqlgentest.AddGqlSchemaFile(t, filepath.Join(testdataDir, "apis", "graphql", "hello", "*.graphqls"))
	gqlgentest.AddGqlSchema("schema.graphqls", `
directive @grpc(service: String!, rpc: String!) on FIELD_DEFINITION
directive @proto(fullName: String!, package: String!, name: String!, goPackage: String!, goName: String!) on OBJECT | INPUT_OBJECT | ENUM
directive @protoField(name: String!, type: String!, goName: String!, goTypeName: String!, goTypePackage: String) on FIELD_DEFINITION | INPUT_FIELD_DEFINITION
scalar DateTime`)
	gqlgentest.AddGqlSchema("hello.graphqls", `
input CreateHelloInput {
  message: String!
}

type CreateHelloPayload {
  hello: Hello!
}

type HellosByUserConnection {
  totalCount: Int!
  edges: [HelloByUserEdge!]!
  nodes: [Hello!]!
  pageInfo: HellosByUserConnectionPageInfo!
}

type HelloByUserEdge {
  node: Hello!
  cursor: String!
}

type HellosByUserConnectionPageInfo {
  endCursor: String!
  hasNextPage: Boolean!
}

extend type Query {
  hellosByUser(userId: Int!): HellosByUserConnection!
}

extend type Mutation {
  createHello(input: CreateHelloInput): CreateHelloPayload!
}`)
	gqlgentest.AddGoModReplace("github.com/izumin5210/remixer", rootDir)
	gqlgentest.AddGoModReplace("apis/go/hello", filepath.Join(testdataDir, "apis", "go", "hello"))

	gqlgentest.Run(t, func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to generate code: %v", err)
		}

		if entries, err := filepath.Glob("**/*"); err != nil {
			t.Errorf("failed to search files: %v", err)
		} else if got, want := len(entries), 6; got != want {
			t.Errorf("found %d files, want %d", got, want)
		}

		gqlgentest.SnapshotFile(t,
			"model/models_gen.go",
			"model/protomodels_gen.go",
			"resolver/resolver.go",
			"resolver/hello.resolvers.go",
			"resolver/hello.resolvers.proto.go",
		)
	})
}

func getModuleRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), ".."))
}

func rewriteFile(t *testing.T, filename string, f func(string, func(string))) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		t.Fatalf("failed to open file for writing: %v", err)
	}

	w := bufio.NewWriter(file)
	f(string(data), func(out string) {
		_, err := w.WriteString(out)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	})

	err = w.Flush()
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

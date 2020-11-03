// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"

	"github.com/izumin5210/remixer/gqlgenplugin"
	"github.com/izumin5210/remixer/gqlgenplugin/protomodelgen"
	"github.com/izumin5210/remixer/gqlgenplugin/protoresolvergen"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	err = api.Generate(cfg,
		gqlgenplugin.AddPluginBefore(protomodelgen.New(), "modelgen"),
		gqlgenplugin.AddPluginBefore(protoresolvergen.New(), "resolvergen"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

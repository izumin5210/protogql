// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"

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
		PrependPlugin(protomodelgen.New()),
		api.AddPlugin(protoresolvergen.New()),
		RemovePlugin("resolvergen"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

func PrependPlugin(p plugin.Plugin) api.Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		*plugins = append([]plugin.Plugin{p}, *plugins...)
	}
}

func RemovePlugin(name string) api.Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		newPlugins := make([]plugin.Plugin, 0, len(*plugins))
		for _, p := range *plugins {
			if p.Name() != name {
				newPlugins = append(newPlugins, p)
			}
		}
		*plugins = newPlugins
	}
}

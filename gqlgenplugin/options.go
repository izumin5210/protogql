package gqlgenplugin

import (
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
)

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

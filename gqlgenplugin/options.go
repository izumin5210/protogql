package gqlgenplugin

import (
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
)

func AddPluginBefore(newPlugin plugin.Plugin, name string) api.Option {
	return func(cfg *config.Config, plugins *[]plugin.Plugin) {
		for i := 0; i < len(*plugins); i++ {
			if p := (*plugins)[i]; p.Name() == name {
				*plugins = append((*plugins)[:i], append([]plugin.Plugin{newPlugin}, (*plugins)[i:]...)...)
				break
			}
		}
	}
}

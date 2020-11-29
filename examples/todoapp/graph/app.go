package graph

import (
	"todoapp/graph/loader"
	"todoapp/graph/resolver"
)

type App struct {
	Resolver *resolver.Resolver
	Loaders  *loader.Loaders
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"testapp/graph"
)

// Hello returns graph.HelloResolver implementation.
func (r *Resolver) Hello() graph.HelloResolver { return &helloResolver{r} }

type helloResolver struct{ *Resolver }


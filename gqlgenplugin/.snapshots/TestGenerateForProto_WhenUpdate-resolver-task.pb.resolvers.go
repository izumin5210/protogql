package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"testapp/graph"
)

// Task returns graph.TaskResolver implementation.
func (r *Resolver) Task() graph.TaskResolver { return &taskResolver{r} }

type taskResolver struct{ *Resolver }


// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package resolver

import (
	"testapp/graph"
)

// Task returns graph.TaskResolver implementation.
func (r *Resolver) Task() graph.TaskResolver { return &taskResolver{r} }

type taskResolver struct{ *Resolver }

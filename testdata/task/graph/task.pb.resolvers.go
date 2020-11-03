package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"task/graph/generated"
)

// Task returns generated.TaskResolver implementation.
func (r *Resolver) Task() generated.TaskResolver { return &taskResolver{r} }

type taskResolver struct{ *Resolver }

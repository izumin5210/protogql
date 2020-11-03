package resolver

import (
	"testapp/graph"
)

// Task returns graph.TaskResolver implementation.
func (r *Resolver) Task() graph.TaskResolver { return &taskProtoResolverAdapter{&taskProtoResolver{r}} }

type taskProtoResolver struct{ *Resolver }


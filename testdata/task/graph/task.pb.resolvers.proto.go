package graph

import (
	"testdata/task/graph/generated"
)

// Task returns generated.TaskResolver implementation.
func (r *Resolver) Task() generated.TaskResolver {
	return &taskProtoResolverAdapter{&taskProtoResolver{r}}
}

type taskProtoResolver struct{ *Resolver }

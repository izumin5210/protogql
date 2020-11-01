package resolver

import (
	task_pb "apis/go/task"
	"context"
	"fmt"
	"testapp/graph"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }


package resolver

import (
	task_pb "apis/go/task"
	"context"
	"fmt"
	"testapp"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns testapp.QueryResolver implementation.
func (r *Resolver) Query() testapp.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }


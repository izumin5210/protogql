package resolver

import (
	task_pb "apis/go/task"
	"context"
	"fmt"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

type queryProtoResolver struct{ *Resolver }


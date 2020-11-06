package resolver

import (
	task_pb "apis/go/task"
	"context"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


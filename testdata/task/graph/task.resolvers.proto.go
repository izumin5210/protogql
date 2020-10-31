package graph

import (
	task_pb "apis/go/task"
	"context"
	"fmt"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *task_pb.Task) ([]*task_pb.User, error) {
	panic(fmt.Errorf("not implemented"))
}

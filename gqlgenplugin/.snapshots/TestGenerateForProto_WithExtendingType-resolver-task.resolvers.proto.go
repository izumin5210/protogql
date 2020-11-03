package resolver

import (
	task_pb "apis/go/task"
	user_pb "apis/go/user"
	"context"
	"fmt"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryProtoResolver) LatestTask(ctx context.Context) (*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *task_pb.Task) ([]*user_pb.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *taskProtoResolver) Author(ctx context.Context, obj *task_pb.Task) (*user_pb.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userProtoResolver) AssignedTasks(ctx context.Context, obj *user_pb.User) ([]*task_pb.Task, error) {
	panic(fmt.Errorf("not implemented"))
}


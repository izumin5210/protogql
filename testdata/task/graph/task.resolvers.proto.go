package graph

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
	"context"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*todo_pb.Task, error) {
	panic("not implemented")
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *todo_pb.Task) ([]*user_pb.User, error) {
	panic("not implemented")
}

func (r *taskProtoResolver) Author(ctx context.Context, obj *todo_pb.Task) (*user_pb.User, error) {
	panic("not implemented")
}

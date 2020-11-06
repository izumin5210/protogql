package resolver

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
	"context"
	_ "net/http/pprof"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*todo_pb.Task, error) {
	return []*todo_pb.Task{}, nil
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *todo_pb.Task) ([]*user_pb.User, error) {
	panic("not implemented")
}

func (r *taskProtoResolver) Author(ctx context.Context, obj *todo_pb.Task) (*user_pb.User, error) {
	panic("not implemented")
}

func (r *userProtoResolver) AssignedTasks(ctx context.Context, obj *user_pb.User) ([]*todo_pb.Task, error) {
	panic("not implemented")
}

func (r *userProtoResolver) TodayTasks(ctx context.Context, obj *user_pb.User) ([]*todo_pb.Task, error) {
	panic("not implemented")
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryProtoResolver) LatestTask(ctx context.Context) (*todo_pb.Task, error) {
	panic("not implemented")
}


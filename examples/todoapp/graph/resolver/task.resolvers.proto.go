package resolver

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
	"context"
	"todoapp/graph/model"
)

func (r *mutationProtoResolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.CreateTaskPayload_Proto, error) {
	req := &todo_pb.CreateTaskRequest{
		Task: model.TaskInputToProto(input.Task),
	}
	resp, err := r.TaskClient.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.CreateTaskPayload_Proto{
		Task: resp.GetTask(),
	}, nil
}

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*todo_pb.Task, error) {
	resp, err := r.TaskClient.ListTasks(ctx, new(todo_pb.ListTasksRequest))
	if err != nil {
		return nil, err
	}
	return resp.GetTasks(), nil
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *todo_pb.Task) ([]*user_pb.User, error) {
	panic("not implemented")
}

func (r *taskProtoResolver) Author(ctx context.Context, obj *todo_pb.Task) (*user_pb.User, error) {
	panic("not implemented")
}

type mutationProtoResolver struct{ *Resolver }
type queryProtoResolver struct{ *Resolver }

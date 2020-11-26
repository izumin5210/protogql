package resolver

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
	"context"
	"todoapp/graph/loader"
	"todoapp/graph/model"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
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
	users, errs := loader.For(ctx).UserByID(ctx).LoadAll(obj.GetAssigneeIds())
	if len(errs) != 0 {
		return nil, multierror.Append(nil, errs...)
	}
	return users, nil
}

func (r *taskProtoResolver) Author(ctx context.Context, obj *todo_pb.Task) (*user_pb.User, error) {
	user, err := loader.For(ctx).UserByID(ctx).Load(obj.GetAuthorId())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

type mutationProtoResolver struct{ *Resolver }
type queryProtoResolver struct{ *Resolver }

package graph

import (
	"context"
	"fmt"
	"testdata/task/api"
	"testdata/task/graph/model"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*api.Task, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *taskProtoResolver) Assignees(ctx context.Context, obj *api.Task) ([]*api.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (a *queryProtoResolverAdapter) Tasks(ctx context.Context) ([]*model.Task, error) {
	resp, err := a.protoResolver.Tasks(ctx)
	if err != nil {
		return nil, err
	}

	return model.TaskListFromRepeatedProto(resp), nil

}
func (a *taskProtoResolverAdapter) Assignees(ctx context.Context, obj *model.Task) ([]*model.User, error) {
	resp, err := a.protoResolver.Assignees(ctx, model.TaskToProto(obj))
	if err != nil {
		return nil, err
	}

	return model.UserListFromRepeatedProto(resp), nil

}

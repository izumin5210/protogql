package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"todo/graph/model"
)

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&queryProtoResolver{Resolver: r.Resolver}).Tasks(ctx)
	if err != nil {
		return nil, err
	}
	return model.TaskListFromRepeatedProto(resp), nil
}

func (r *taskResolver) Assignees(ctx context.Context, obj *model.Task) ([]*model.User, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&taskProtoResolver{Resolver: r.Resolver}).Assignees(ctx, model.TaskToProto(obj))
	if err != nil {
		return nil, err
	}
	return model.UserListFromRepeatedProto(resp), nil
}

func (r *taskResolver) Author(ctx context.Context, obj *model.Task) (*model.User, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&taskProtoResolver{Resolver: r.Resolver}).Author(ctx, model.TaskToProto(obj))
	if err != nil {
		return nil, err
	}
	return model.UserFromProto(resp), nil
}

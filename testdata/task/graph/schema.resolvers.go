package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"testdata/task/graph/generated"
	"testdata/task/graph/model"
)

func (r *mutationResolver) Nop(ctx context.Context) (*bool, error) {
	return r.ProtoResolver.Mutation().Nop(ctx)
}

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	out, err := r.ProtoResolver.Query().Tasks(ctx)
	if err != nil {
		return nil, err
	}
	return model.TaskListFromRepeatedProto(out), nil
}

func (r *taskResolver) Assignees(ctx context.Context, obj *model.Task) ([]*model.User, error) {
	out, err := r.ProtoResolver.Task().Assignees(ctx, obj.Proto())
	if err != nil {
		return nil, err
	}
	return model.UserListFromRepeatedProto(out), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"testapp/graph"
	"testapp/model"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input *model.CreateTaskInput) (*model.CreateTaskPayload, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&mutationProtoResolver{Resolver: r.Resolver}).CreateTask(ctx, input)
	if err != nil {
		return nil, err
	}
	return model.CreateTaskPayloadFromProto(resp), nil
}

func (r *queryResolver) TasksByUser(ctx context.Context, userID int) (*model.TasksByUserConnection, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&queryProtoResolver{Resolver: r.Resolver}).TasksByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.TasksByUserConnectionFromProto(resp), nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"task/graph/generated"
)

func (r *mutationResolver) Nop(ctx context.Context) (*bool, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&mutationProtoResolver{Resolver: r.Resolver}).Nop(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *queryResolver) Nop(ctx context.Context) (*bool, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&queryProtoResolver{Resolver: r.Resolver}).Nop(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
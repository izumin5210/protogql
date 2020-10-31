package graph

import (
	"context"
	"fmt"
	"task/graph/generated"
)

func (r *mutationProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationProtoResolverAdapter{&mutationProtoResolver{r}}
}

type mutationProtoResolver struct{ *Resolver }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }

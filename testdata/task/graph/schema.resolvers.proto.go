package graph

import (
	"context"
)

func (r *mutationProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic("not implemented")
}

func (r *queryProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic("not implemented")
}

type mutationProtoResolver struct{ *Resolver }
type queryProtoResolver struct{ *Resolver }

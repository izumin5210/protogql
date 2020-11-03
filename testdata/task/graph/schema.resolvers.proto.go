package graph

import (
	"context"
	"fmt"
)

func (r *mutationProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

type mutationProtoResolver struct{ *Resolver }
type queryProtoResolver struct{ *Resolver }

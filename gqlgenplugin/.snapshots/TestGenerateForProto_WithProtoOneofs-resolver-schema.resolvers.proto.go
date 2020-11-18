package resolver

import (
	oneof_pb "apis/go/oneof"
	"context"
)

func (r *queryProtoResolver) Entries(ctx context.Context) ([]*oneof_pb.Entry, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


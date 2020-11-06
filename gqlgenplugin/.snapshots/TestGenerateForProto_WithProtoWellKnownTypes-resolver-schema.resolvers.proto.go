package resolver

import (
	wktypes_pb "apis/go/wktypes"
	"context"
)

func (r *queryProtoResolver) Hello(ctx context.Context) ([]*wktypes_pb.Hello, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


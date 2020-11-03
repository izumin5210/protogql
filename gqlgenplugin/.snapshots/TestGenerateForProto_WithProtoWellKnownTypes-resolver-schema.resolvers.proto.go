package resolver

import (
	wktypes_pb "apis/go/wktypes"
	"context"
	"fmt"
)

func (r *queryProtoResolver) Hello(ctx context.Context) ([]*wktypes_pb.Hello, error) {
	panic(fmt.Errorf("not implemented"))
}

type queryProtoResolver struct{ *Resolver }


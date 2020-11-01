package resolver

import (
	wktypes_pb "apis/go/wktypes"
	"context"
	"fmt"
	"testapp/graph"
)

func (r *queryProtoResolver) Hello(ctx context.Context) ([]*wktypes_pb.Hello, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }


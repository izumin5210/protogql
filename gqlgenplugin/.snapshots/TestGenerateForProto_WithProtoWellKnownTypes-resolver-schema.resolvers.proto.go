package resolver

import (
	wktypes_pb "apis/go/wktypes"
	"context"
	"fmt"
	"testapp"
)

func (r *queryProtoResolver) Hello(ctx context.Context) ([]*wktypes_pb.Hello, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns testapp.QueryResolver implementation.
func (r *Resolver) Query() testapp.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }


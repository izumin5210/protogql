// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package resolver

import (
	"context"
	"testapp/graph"
	"testapp/model"
)

func (r *queryResolver) Hello(ctx context.Context) ([]*model.Hello, error) {
	// This function body is generated by github.com/izumin5210/remixer. DO NOT EDIT.

	resp, err := (&queryProtoResolver{Resolver: r.Resolver}).Hello(ctx)
	if err != nil {
		return nil, err
	}
	return model.HelloListFromRepeatedProto(resp), nil
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

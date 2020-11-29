package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"testapp/graph"
	"testapp/model"
)

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	// This function body is generated by github.com/izumin5210/protogql. DO NOT EDIT.

	resp, err := (&queryProtoResolver{Resolver: r.Resolver}).Hello(ctx)
	if err != nil {
		return nil, err
	}
	return model.HelloFromProto(resp), nil
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }


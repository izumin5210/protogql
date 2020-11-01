package resolver

import (
	"apis/go/task"
	"context"
	"fmt"
	"testapp"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*task.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns testapp.QueryResolver implementation.
func (r *Resolver) Query() testapp.QueryResolver {
	return &queryProtoResolverAdapter{&queryProtoResolver{r}}
}

type queryProtoResolver struct{ *Resolver }


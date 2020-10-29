package graph

import (
	"context"
	"fmt"
	"testdata/task/api"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*api.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *api.Task) ([]*api.User, error) {
	panic(fmt.Errorf("not implemented"))
}

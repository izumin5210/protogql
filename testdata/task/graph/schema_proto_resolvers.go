package graph

import (
	"context"
	"fmt"

	"testdata/task/api"
)

func (r *mutationProtoResolver) Nop(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*api.Task, error) {
	return []*api.Task{
		{Id: 123, Title: "Foo", AssigneeIds: []uint64{2, 4}},
		{Id: 124, Title: "Bar", AssigneeIds: []uint64{1, 4, 6}},
	}, nil
}

func (r *taskProtoResolver) Assignees(ctx context.Context, obj *api.Task) ([]*api.User, error) {
	return nil, nil
}

func (r *ProtoResolver) Mutation() *mutationProtoResolver { return &mutationProtoResolver{r} }
func (r *ProtoResolver) Query() *queryProtoResolver       { return &queryProtoResolver{r} }

type mutationProtoResolver struct{ *ProtoResolver }
type queryProtoResolver struct{ *ProtoResolver }

func (r *ProtoResolver) Task() *taskProtoResolver { return &taskProtoResolver{r} }

type taskProtoResolver struct{ *ProtoResolver }

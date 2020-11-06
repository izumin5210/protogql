package resolver

import (
	todo_pb "apis/go/todo"
	"context"
)

func (r *queryProtoResolver) Tasks(ctx context.Context) ([]*todo_pb.Task, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


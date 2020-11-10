package resolver

import (
	"context"
	"testapp/model"
)

func (r *mutationProtoResolver) CreateTask(ctx context.Context, input *model.CreateTaskInput) (*model.CreateTaskPayload_Proto, error) {
	panic("not implemented")
}

type mutationProtoResolver struct{ *Resolver }


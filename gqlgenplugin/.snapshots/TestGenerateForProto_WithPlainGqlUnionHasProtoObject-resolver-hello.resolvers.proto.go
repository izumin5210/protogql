package resolver

import (
	"context"
	"testapp/model"
)

func (r *mutationProtoResolver) CreateHello(ctx context.Context, input *model.CreateHelloInput) (*model.CreateHelloPayload_Proto, error) {
	panic("not implemented")
}

type mutationProtoResolver struct{ *Resolver }


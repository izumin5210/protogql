package resolver

import (
	user_pb "apis/go/user"
	"context"
)

func (r *queryProtoResolver) CurrentUser(ctx context.Context) (*user_pb.User, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


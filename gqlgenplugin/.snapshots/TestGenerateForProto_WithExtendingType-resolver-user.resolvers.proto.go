package resolver

import (
	user_pb "apis/go/user"
	"context"
	"fmt"
)

func (r *queryProtoResolver) CurrentUser(ctx context.Context) (*user_pb.User, error) {
	panic(fmt.Errorf("not implemented"))
}

type queryProtoResolver struct{ *Resolver }


package resolver

import (
	hello_pb "apis/go/hello"
	"context"
)

func (r *helloProtoResolver) User(ctx context.Context, obj *hello_pb.Hello) (*hello_pb.User, error) {
	panic("not implemented")
}

func (r *queryProtoResolver) Hello(ctx context.Context) (*hello_pb.Hello, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }


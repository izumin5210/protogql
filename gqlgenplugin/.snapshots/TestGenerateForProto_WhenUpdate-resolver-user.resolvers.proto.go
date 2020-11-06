package resolver

import (
	user_pb "apis/go/user"
	"context"
)

func (r *queryProtoResolver) CurrentUser(ctx context.Context) (*user_pb.User, error) {
	panic("not implemented")
}

type queryProtoResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const (
	TestConstant = 1
)

var (
	TestVariable = 1
)

type TestStruct struct {
	Foo string
}

func TestFunction() string { return "Test" }


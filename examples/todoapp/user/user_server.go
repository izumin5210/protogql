package user

import (
	"context"

	user_pb "apis/go/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Users = map[uint64]*user_pb.User{}

func NewUserServiceServer() user_pb.UserServiceServer {
	return new(userServer)
}

type userServer struct {
	user_pb.UserServiceServer
}

func (s *userServer) BatchGetUsers(ctx context.Context, in *user_pb.BatchGetUsersRequest) (*user_pb.BatchGetUsersResponse, error) {
	resp := &user_pb.BatchGetUsersResponse{Users: make([]*user_pb.User, len(in.GetUserIds()))}

	for i, id := range in.GetUserIds() {
		user, ok := Users[id]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "user(id=%d) not found", id)
		}
		resp.Users[i] = user
	}

	return resp, nil
}

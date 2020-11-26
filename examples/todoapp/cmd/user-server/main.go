package main

import (
	"context"
	"log"
	"net"
	"os"

	user_pb "apis/go/user"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("USER_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user_pb.RegisterUserServiceServer(s, NewUserServiceServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewUserServiceServer() user_pb.UserServiceServer {
	return new(userServer)
}

type userServer struct {
	user_pb.UserServiceServer
}

func (s *userServer) BatchGetUsers(ctx context.Context, in *user_pb.BatchGetUsersRequest) (*user_pb.BatchGetUsersResponse, error) {
	resp := new(user_pb.BatchGetUsersResponse)
	return resp, nil
}

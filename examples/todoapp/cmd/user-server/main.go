package main

import (
	"log"
	"net"
	"os"

	user_pb "apis/go/user"

	"google.golang.org/grpc"

	"todoapp/user"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("USER_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user_pb.RegisterUserServiceServer(s, user.NewUserServiceServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

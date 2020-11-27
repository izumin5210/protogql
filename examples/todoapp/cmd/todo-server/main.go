package main

import (
	"log"
	"net"
	"os"

	todo_pb "apis/go/todo"

	"google.golang.org/grpc"

	"todoapp/todo"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("TASK_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	todo_pb.RegisterTaskServiceServer(s, todo.NewTaskServiceServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

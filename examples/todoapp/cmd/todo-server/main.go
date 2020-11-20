package main

import (
	"context"
	"log"
	"net"
	"os"

	todo_pb "apis/go/todo"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("TASK_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	todo_pb.RegisterTaskServiceServer(s, NewTaskServiceServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewTaskServiceServer() todo_pb.TaskServiceServer {
	return new(taskServer)
}

type taskServer struct {
	todo_pb.TaskServiceServer
}

func (s *taskServer) ListTasks(context.Context, *todo_pb.ListTasksRequest) (*todo_pb.ListTasksResponse, error) {
	resp := &todo_pb.ListTasksResponse{}
	return resp, nil
}

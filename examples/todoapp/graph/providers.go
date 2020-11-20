package graph

import (
	"os"

	todo_pb "apis/go/todo"

	"github.com/google/wire"
	"google.golang.org/grpc"

	"todoapp/graph/resolver"
)

func provideTaskClient() (todo_pb.TaskServiceClient, func(), error) {
	conn, err := grpc.Dial("localhost:"+os.Getenv("TASK_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}
	return todo_pb.NewTaskServiceClient(conn), func() { conn.Close() }, nil
}

var appSet = wire.NewSet(
	wire.Struct(new(App), "*"),
	wire.Struct(new(resolver.Resolver), "*"),
	provideTaskClient,
)

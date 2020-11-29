package graph

import (
	"os"

	todo_pb "apis/go/todo"
	user_pb "apis/go/user"

	"github.com/google/wire"
	"google.golang.org/grpc"

	"todoapp/graph/loader"
	"todoapp/graph/resolver"
)

func provideTaskClient() (todo_pb.TaskServiceClient, func(), error) {
	conn, err := grpc.Dial("localhost:"+os.Getenv("TASK_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}
	return todo_pb.NewTaskServiceClient(conn), func() { conn.Close() }, nil
}

func provideUserClient() (user_pb.UserServiceClient, func(), error) {
	conn, err := grpc.Dial("localhost:"+os.Getenv("USER_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}
	return user_pb.NewUserServiceClient(conn), func() { conn.Close() }, nil
}

var appSet = wire.NewSet(
	wire.Struct(new(App), "*"),
	wire.Struct(new(resolver.Resolver), "*"),
	wire.Struct(new(loader.Loaders), "*"),
	provideTaskClient,
	provideUserClient,
)

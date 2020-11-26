package resolver

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TaskClient todo_pb.TaskServiceClient
	UserClient user_pb.UserServiceClient
}

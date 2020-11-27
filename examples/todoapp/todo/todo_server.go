package todo

import (
	"context"
	"sort"

	todo_pb "apis/go/todo"
)

var Tasks = map[uint64]*todo_pb.Task{}

func NewTaskServiceServer() todo_pb.TaskServiceServer {
	return new(taskServer)
}

type taskServer struct {
	todo_pb.TaskServiceServer
}

func (s *taskServer) ListTasks(ctx context.Context, in *todo_pb.ListTasksRequest) (*todo_pb.ListTasksResponse, error) {
	resp := &todo_pb.ListTasksResponse{}

	ids := make([]uint64, 0, len(Tasks))
	for id := range Tasks {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		resp.Tasks = append(resp.Tasks, Tasks[id])
	}

	return resp, nil
}
